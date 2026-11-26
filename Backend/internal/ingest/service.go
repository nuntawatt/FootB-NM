// Package ingest persists football fixture data fetched from an external
// provider into the local PostgreSQL database.
//
// Entry point: Service.SyncFixtures — iterates over all leagues that have an
// external_api_id, fetches their fixtures, and upserts teams and matches
// inside a single database transaction per league.
package ingest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"football/internal/models"
	"football/internal/provider"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ----------------------------------------------------------------------------
// Service
// ----------------------------------------------------------------------------

// Service orchestrates fixture ingestion. Create one with New.
type Service struct {
	db  *gorm.DB
	p   provider.Provider
	log *slog.Logger
}

// New creates a Service. If logger is nil, slog.Default() is used.
func New(db *gorm.DB, p provider.Provider, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{db: db, p: p, log: logger}
}

// ----------------------------------------------------------------------------
// Result
// ----------------------------------------------------------------------------

// Result summarises the outcome of a sync operation.
type Result struct {
	Leagues        int // number of leagues successfully synced
	TeamsCreated   int
	MatchesCreated int
	MatchesUpdated int
	Skipped        int // fixtures that were skipped due to invalid data
}

// add merges another Result into this one (used to aggregate per-league totals).
func (r *Result) add(other Result) {
	r.Leagues += other.Leagues
	r.TeamsCreated += other.TeamsCreated
	r.MatchesCreated += other.MatchesCreated
	r.MatchesUpdated += other.MatchesUpdated
	r.Skipped += other.Skipped
}

// ----------------------------------------------------------------------------
// Public API
// ----------------------------------------------------------------------------

// SyncFixtures fetches and persists fixtures for every league that has an
// external_api_id. Leagues are processed in sort order (league before cup).
//
// If one league fails, the error is collected and remaining leagues are still
// attempted. The returned error is a joined error of all failures.
func (s *Service) SyncFixtures(ctx context.Context, season int, from, to time.Time) (Result, error) {
	leagues, err := s.loadLinkedLeagues(ctx)
	if err != nil {
		return Result{}, err
	}

	var (
		total      Result
		syncErrors []error
	)

	for _, league := range leagues {
		leagueResult, err := s.syncLeague(ctx, league, season, from, to)
		if err != nil {
			syncErrors = append(syncErrors, err)
			continue
		}

		total.Leagues++
		total.add(leagueResult)

		s.log.Info("league synced",
			"league", league.Name,
			"matches_created", leagueResult.MatchesCreated,
			"matches_updated", leagueResult.MatchesUpdated,
			"teams_created", leagueResult.TeamsCreated,
			"skipped", leagueResult.Skipped,
		)
	}

	return total, errors.Join(syncErrors...)
}

// ApplyFixtures persists a pre-fetched slice of fixtures for a single league.
// This is useful when the caller already holds the provider.Fixture slice
// (e.g. from a live-fixtures poll) and wants to skip the HTTP fetch step.
func (s *Service) ApplyFixtures(ctx context.Context, league models.League, fixtures []provider.Fixture) (Result, error) {
	return s.persistFixtures(ctx, league, fixtures)
}

// ----------------------------------------------------------------------------
// Internal helpers
// ----------------------------------------------------------------------------

// loadLinkedLeagues returns all leagues that have an external_api_id set,
// ordered so that "league" type comes before "cup" type.
func (s *Service) loadLinkedLeagues(ctx context.Context) ([]models.League, error) {
	var leagues []models.League

	err := s.db.WithContext(ctx).
		Where("external_api_id IS NOT NULL").
		Order("CASE WHEN type = 'league' THEN 0 ELSE 1 END, sort_order, name").
		Find(&leagues).Error

	if err != nil {
		return nil, fmt.Errorf("load linked leagues: %w", err)
	}
	if len(leagues) == 0 {
		return nil, errors.New("no linked leagues found; run sync link first")
	}
	return leagues, nil
}

// syncLeague fetches fixtures for one league from the provider and persists them.
func (s *Service) syncLeague(ctx context.Context, league models.League, season int, from, to time.Time) (Result, error) {
	fixtures, err := s.p.Fixtures(ctx, provider.FixtureQuery{
		LeagueExternalID: *league.ExternalAPIID,
		Season:           season,
		From:             from,
		To:               to,
	})
	if err != nil {
		return Result{}, fmt.Errorf("fetch fixtures for league %q: %w", league.Name, err)
	}

	result, err := s.persistFixtures(ctx, league, fixtures)
	if err != nil {
		return Result{}, fmt.Errorf("persist fixtures for league %q: %w", league.Name, err)
	}
	return result, nil
}

// persistFixtures upserts all fixtures for a league inside a single transaction.
func (s *Service) persistFixtures(ctx context.Context, league models.League, fixtures []provider.Fixture) (Result, error) {
	var result Result

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, fixture := range fixtures {
			if !isUsableFixture(fixture) {
				s.log.Warn("skip invalid fixture",
					"league", league.Name,
					"fixture_id", fixture.ExternalID,
				)
				result.Skipped++
				continue
			}

			homeID, created, err := upsertTeam(tx, league.ID, fixture.Home)
			if err != nil {
				return fmt.Errorf("upsert home team: %w", err)
			}
			if created {
				result.TeamsCreated++
			}

			awayID, created, err := upsertTeam(tx, league.ID, fixture.Away)
			if err != nil {
				return fmt.Errorf("upsert away team: %w", err)
			}
			if created {
				result.TeamsCreated++
			}

			created, err = upsertMatch(tx, league.ID, homeID, awayID, fixture)
			if err != nil {
				return fmt.Errorf("upsert match %s: %w", fixture.ExternalID, err)
			}
			if created {
				result.MatchesCreated++
			} else {
				result.MatchesUpdated++
			}
		}
		return nil
	})

	if err != nil {
		return Result{}, fmt.Errorf("transaction failed: %w", err)
	}
	return result, nil
}

// isUsableFixture returns true when the fixture has all required fields.
// Fixtures with zero/missing IDs, identical home/away teams, or a missing
// kick-off time are skipped to prevent bad data entering the database.
func isUsableFixture(f provider.Fixture) bool {
	isBlank := func(id string) bool { return id == "" || id == "0" }

	return !isBlank(f.ExternalID) &&
		!isBlank(f.Home.ExternalID) &&
		!isBlank(f.Away.ExternalID) &&
		f.Home.ExternalID != f.Away.ExternalID &&
		!f.KickOff.IsZero()
}

// upsertTeam finds an existing team by external_api_id or creates a new one.
// If the team already exists, only name and logo_url are updated when they change.
// LeagueID is intentionally NOT updated: a team may appear in multiple
// competitions without changing its primary league association.
//
// Returns the team's UUID, whether it was newly created, and any error.
func upsertTeam(tx *gorm.DB, leagueID uuid.UUID, data provider.Team) (uuid.UUID, bool, error) {
	var team models.Team

	err := tx.Where("external_api_id = ?", data.ExternalID).First(&team).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create a new team record.
		externalID := data.ExternalID
		team = models.Team{
			LeagueID:      leagueID,
			Name:          data.Name,
			LogoURL:       data.LogoURL,
			ExternalAPIID: &externalID,
		}
		if err := tx.Create(&team).Error; err != nil {
			return uuid.Nil, false, fmt.Errorf("create team: %w", err)
		}
		return team.ID, true, nil
	}

	if err != nil {
		return uuid.Nil, false, fmt.Errorf("find team: %w", err)
	}

	// Update only the fields that have actually changed.
	updates := make(map[string]any)
	if data.Name != "" && data.Name != team.Name {
		updates["name"] = data.Name
	}
	if data.LogoURL != "" && data.LogoURL != team.LogoURL {
		updates["logo_url"] = data.LogoURL
	}
	if len(updates) > 0 {
		if err := tx.Model(&team).Updates(updates).Error; err != nil {
			return uuid.Nil, false, fmt.Errorf("update team: %w", err)
		}
	}

	return team.ID, false, nil
}

// upsertMatch finds an existing match by external_api_id or creates a new one.
// Every field is updated on subsequent syncs so that live score changes and
// status transitions are reflected in the database.
//
// Returns whether the record was newly created and any error.
func upsertMatch(tx *gorm.DB, leagueID, homeID, awayID uuid.UUID, f provider.Fixture) (bool, error) {
	var match models.Match

	err := tx.Where("external_api_id = ?", f.ExternalID).First(&match).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Default to "scheduled" when the provider status is missing or unknown.
		status := f.Status
		if status == "" {
			status = provider.StatusScheduled
		}

		externalID := f.ExternalID
		match = models.Match{
			LeagueID:         leagueID,
			HomeTeamID:       homeID,
			AwayTeamID:       awayID,
			MatchDate:        f.KickOff,
			HomeScore:        toInt16(f.HomeScore),
			AwayScore:        toInt16(f.AwayScore),
			HomePenaltyScore: toInt16(f.HomePenalty),
			AwayPenaltyScore: toInt16(f.AwayPenalty),
			Minute:           toInt16(f.Minute),
			Round:            f.Round,
			Status:           status,
			ExternalAPIID:    &externalID,
		}
		if err := tx.Create(&match).Error; err != nil {
			return false, fmt.Errorf("create match: %w", err)
		}
		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf("find match: %w", err)
	}

	// Update all mutable fields; status is only updated when non-empty.
	updates := map[string]any{
		"league_id":          leagueID,
		"home_team_id":       homeID,
		"away_team_id":       awayID,
		"match_date":         f.KickOff,
		"home_score":         toInt16(f.HomeScore),
		"away_score":         toInt16(f.AwayScore),
		"home_penalty_score": toInt16(f.HomePenalty),
		"away_penalty_score": toInt16(f.AwayPenalty),
		"minute":             toInt16(f.Minute),
		"round":              f.Round,
	}
	if f.Status != "" {
		updates["status"] = f.Status
	}

	if err := tx.Model(&match).Updates(updates).Error; err != nil {
		return false, fmt.Errorf("update match: %w", err)
	}
	return false, nil
}

// toInt16 converts a *int to *int16.
// Returns nil when the input is nil, preserving the "not available" semantic.
func toInt16(v *int) *int16 {
	if v == nil {
		return nil
	}
	n := int16(*v)
	return &n
}
