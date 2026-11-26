
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

type Service struct {
	db  *gorm.DB
	p   provider.Provider
	log *slog.Logger
}

type Result struct {
	Leagues        int
	TeamsCreated   int
	MatchesCreated int
	MatchesUpdated int
	Skipped        int
}

func New(
	db *gorm.DB,
	p provider.Provider,
	logger *slog.Logger,
) *Service {
	if logger == nil {
		logger = slog.Default()
	}

	return &Service{
		db:  db,
		p:   p,
		log: logger,
	}
}

// SyncFixtures ดึงข้อมูลการแข่งขันจากทุกลีกที่ผูก external_api_id ไว้
func (s *Service) SyncFixtures(
	ctx context.Context,
	season int,
	from, to time.Time,
) (Result, error) {
	var result Result

	var leagues []models.League

	err := s.db.WithContext(ctx).
		Where("external_api_id IS NOT NULL").
		Order("CASE WHEN type = 'league' THEN 0 ELSE 1 END, sort_order, name").
		Find(&leagues).Error

	if err != nil {
		return result, fmt.Errorf("find leagues: %w", err)
	}

	if len(leagues) == 0 {
		return result, errors.New(
			"no linked leagues found (use sync link first)",
		)
	}

	var syncErrors []error

	for _, league := range leagues {
		leagueResult, err := s.syncLeague(
			ctx,
			league,
			season,
			from,
			to,
		)

		if err != nil {
			syncErrors = append(syncErrors, err)
			continue
		}

		result.Leagues++
		result.add(leagueResult)

		s.log.Info(
			"league synced",
			"league", league.Name,
			"matches_created", leagueResult.MatchesCreated,
			"matches_updated", leagueResult.MatchesUpdated,
			"teams_created", leagueResult.TeamsCreated,
			"skipped", leagueResult.Skipped,
		)
	}

	return result, errors.Join(syncErrors...)
}

// syncLeague ดึง Fixtures และบันทึกข้อมูลของลีกหนึ่งลีก
func (s *Service) syncLeague(
	ctx context.Context,
	league models.League,
	season int,
	from, to time.Time,
) (Result, error) {
	var result Result

	fixtures, err := s.p.Fixtures(
		ctx,
		provider.FixtureQuery{
			LeagueExternalID: *league.ExternalAPIID,
			Season:           season,
			From:             from,
			To:               to,
		},
	)

	if err != nil {
		return result, fmt.Errorf(
			"fetch fixtures for league %q: %w",
			league.Name,
			err,
		)
	}

	result, err = s.apply(ctx, league, fixtures)
	if err != nil {
		return Result{}, fmt.Errorf(
			"apply fixtures for league %q: %w",
			league.Name,
			err,
		)
	}

	return result, nil
}

// ApplyFixtures บันทึก Fixtures ที่ได้รับจาก Provider
func (s *Service) ApplyFixtures(
	ctx context.Context,
	league models.League,
	fixtures []provider.Fixture,
) (Result, error) {
	return s.apply(ctx, league, fixtures)
}

// apply บันทึกข้อมูลการแข่งขันภายใน Transaction
func (s *Service) apply(
	ctx context.Context,
	league models.League,
	fixtures []provider.Fixture,
) (Result, error) {
	var result Result

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var txResult Result

		for _, fixture := range fixtures {
			if !usable(fixture) {
				s.log.Warn(
					"skip invalid fixture",
					"league", league.Name,
					"fixture", fixture.ExternalID,
				)

				txResult.Skipped++
				continue
			}

			homeID, created, err := upsertTeam(
				tx,
				league.ID,
				fixture.Home,
			)
			if err != nil {
				return fmt.Errorf("upsert home team: %w", err)
			}

			if created {
				txResult.TeamsCreated++
			}

			awayID, created, err := upsertTeam(
				tx,
				league.ID,
				fixture.Away,
			)
			if err != nil {
				return fmt.Errorf("upsert away team: %w", err)
			}

			if created {
				txResult.TeamsCreated++
			}

			created, err = upsertMatch(
				tx,
				league.ID,
				homeID,
				awayID,
				fixture,
			)
			if err != nil {
				return fmt.Errorf(
					"upsert match %s: %w",
					fixture.ExternalID,
					err,
				)
			}

			if created {
				txResult.MatchesCreated++
			} else {
				txResult.MatchesUpdated++
			}
		}

		result = txResult
		return nil
	})

	if err != nil {
		return Result{}, fmt.Errorf(
			"transaction failed: %w",
			err,
		)
	}

	return result, nil
}

// add รวมผลลัพธ์จากการ Sync
func (r *Result) add(other Result) {
	r.Leagues += other.Leagues
	r.TeamsCreated += other.TeamsCreated
	r.MatchesCreated += other.MatchesCreated
	r.MatchesUpdated += other.MatchesUpdated
	r.Skipped += other.Skipped
}

// usable ตรวจสอบว่า Fixture มีข้อมูลที่จำเป็นครบหรือไม่
func usable(f provider.Fixture) bool {
	invalidID := func(id string) bool {
		return id == "" || id == "0"
	}

	return !invalidID(f.ExternalID) &&
		!invalidID(f.Home.ExternalID) &&
		!invalidID(f.Away.ExternalID) &&
		f.Home.ExternalID != f.Away.ExternalID &&
		!f.KickOff.IsZero()
}

// upsertTeam ค้นหาและสร้างหรืออัปเดตทีม
func upsertTeam(
	tx *gorm.DB,
	leagueID uuid.UUID,
	teamData provider.Team,
) (uuid.UUID, bool, error) {
	var team models.Team

	err := tx.
		Where("external_api_id = ?", teamData.ExternalID).
		First(&team).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		externalID := teamData.ExternalID

		team = models.Team{
			LeagueID:      leagueID,
			Name:          teamData.Name,
			LogoURL:       teamData.LogoURL,
			ExternalAPIID: &externalID,
		}

		if err := tx.Create(&team).Error; err != nil {
			return uuid.Nil, false, fmt.Errorf(
				"create team: %w",
				err,
			)
		}

		return team.ID, true, nil
	}

	if err != nil {
		return uuid.Nil, false, fmt.Errorf(
			"find team: %w",
			err,
		)
	}

	// ไม่เปลี่ยน LeagueID ของทีมเดิม
	updates := map[string]any{}

	if teamData.Name != "" && teamData.Name != team.Name {
		updates["name"] = teamData.Name
	}

	if teamData.LogoURL != "" && teamData.LogoURL != team.LogoURL {
		updates["logo_url"] = teamData.LogoURL
	}

	if len(updates) > 0 {
		if err := tx.Model(&team).Updates(updates).Error; err != nil {
			return uuid.Nil, false, fmt.Errorf(
				"update team: %w",
				err,
			)
		}
	}

	return team.ID, false, nil
}

// upsertMatch ค้นหาและสร้างหรืออัปเดตการแข่งขัน
func upsertMatch(
	tx *gorm.DB,
	leagueID uuid.UUID,
	homeID uuid.UUID,
	awayID uuid.UUID,
	fixture provider.Fixture,
) (bool, error) {
	var match models.Match

	err := tx.
		Where("external_api_id = ?", fixture.ExternalID).
		First(&match).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		externalID := fixture.ExternalID

		status := fixture.Status
		if status == "" {
			status = provider.StatusScheduled
		}

		match = models.Match{
			LeagueID:          leagueID,
			HomeTeamID:        homeID,
			AwayTeamID:        awayID,
			MatchDate:         fixture.KickOff,
			HomeScore:         i16(fixture.HomeScore),
			AwayScore:         i16(fixture.AwayScore),
			HomePenaltyScore:  i16(fixture.HomePenalty),
			AwayPenaltyScore:  i16(fixture.AwayPenalty),
			Minute:            i16(fixture.Minute),
			Round:             fixture.Round,
			Status:            status,
			ExternalAPIID:     &externalID,
		}

		if err := tx.Create(&match).Error; err != nil {
			return false, fmt.Errorf(
				"create match: %w",
				err,
			)
		}

		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf(
			"find match: %w",
			err,
		)
	}

	updates := map[string]any{
		"league_id":          leagueID,
		"home_team_id":       homeID,
		"away_team_id":       awayID,
		"match_date":         fixture.KickOff,
		"home_score":         i16(fixture.HomeScore),
		"away_score":         i16(fixture.AwayScore),
		"home_penalty_score": i16(fixture.HomePenalty),
		"away_penalty_score": i16(fixture.AwayPenalty),
		"minute":             i16(fixture.Minute),
		"round":              fixture.Round,
	}

	if fixture.Status != "" {
		updates["status"] = fixture.Status
	}

	if err := tx.Model(&match).Updates(updates).Error; err != nil {
		return false, fmt.Errorf(
			"update match: %w",
			err,
		)
	}

	return false, nil
}

// i16 แปลง *int เป็น *int16
func i16(value *int) *int16 {
	if value == nil {
		return nil
	}

	converted := int16(*value)
	return &converted
}
