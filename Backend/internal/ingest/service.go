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

func New(db *gorm.DB, p provider.Provider, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{db: db, p: p, log: logger}
}

type Result struct {
	Leagues        int
	TeamsCreated   int
	MatchesCreated int
	MatchesUpdated int
	Skipped        int
}

func (r *Result) add(other Result) {
	r.Leagues += other.Leagues
	r.TeamsCreated += other.TeamsCreated
	r.MatchesCreated += other.MatchesCreated
	r.MatchesUpdated += other.MatchesUpdated
	r.Skipped += other.Skipped
}

// SyncFixtures ดึงข้อมูลและบันทึก fixtures ของทุกลีกที่มี external_api_id
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

func (s *Service) ApplyFixtures(ctx context.Context, league models.League, fixtures []provider.Fixture) (Result, error) {
	return s.persistFixtures(ctx, league, fixtures)
}

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

// persistFixtures บันทึกข้อมูล fixtures ในระดับ database transaction
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

// isUsableFixture ตรวจสอบความถูกต้องของข้อมูล fixture ก่อนบันทึก
func isUsableFixture(f provider.Fixture) bool {
	isBlank := func(id string) bool { return id == "" || id == "0" }

	return !isBlank(f.ExternalID) &&
		!isBlank(f.Home.ExternalID) &&
		!isBlank(f.Away.ExternalID) &&
		f.Home.ExternalID != f.Away.ExternalID &&
		!f.KickOff.IsZero()
}

// upsertTeam ค้นหาหรือสร้างทีมใหม่ (ไม่เปลี่ยน LeagueID เดิม)
func upsertTeam(tx *gorm.DB, leagueID uuid.UUID, data provider.Team) (uuid.UUID, bool, error) {
	var teams []models.Team
	if err := tx.Where("external_api_id = ?", data.ExternalID).Limit(1).Find(&teams).Error; err != nil {
		return uuid.Nil, false, fmt.Errorf("find team: %w", err)
	}

	if len(teams) == 0 {
		externalID := data.ExternalID
		team := models.Team{
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

	team := teams[0]
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

// upsertMatch ค้นหาและบันทึกหรืออัปเดตผลการแข่งขัน
func upsertMatch(tx *gorm.DB, leagueID, homeID, awayID uuid.UUID, f provider.Fixture) (bool, error) {
	var matches []models.Match
	if err := tx.Where("external_api_id = ?", f.ExternalID).Limit(1).Find(&matches).Error; err != nil {
		return false, fmt.Errorf("find match: %w", err)
	}

	if len(matches) == 0 {
		status := f.Status
		if status == "" {
			status = provider.StatusScheduled
		}

		externalID := f.ExternalID
		match := models.Match{
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

	match := matches[0]
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

func toInt16(v *int) *int16 {
	if v == nil {
		return nil
	}
	n := int16(*v)
	return &n
}

type SquadSyncResult struct {
	TeamsCount   int
	CoachesCount int
	PlayersCount int
}

// SyncSquadsAndCoaches ดึงโค้ชและนักเตะของทุกทีมที่มี external_api_id ลง Database
func (s *Service) SyncSquadsAndCoaches(ctx context.Context) (SquadSyncResult, error) {
	var result SquadSyncResult
	var teams []models.Team

	if err := s.db.WithContext(ctx).Where("external_api_id IS NOT NULL").Find(&teams).Error; err != nil {
		return result, fmt.Errorf("load teams: %w", err)
	}

	var syncErrors []error
	for _, team := range teams {
		if team.ExternalAPIID == nil || *team.ExternalAPIID == "" {
			continue
		}
		result.TeamsCount++
		extID := *team.ExternalAPIID

		// 1. Sync Coach (ถ้ามีอยู่แล้วให้ข้ามเพื่อประหยัดโควตา API)
		var coachCount int64
		s.db.WithContext(ctx).Model(&models.Coach{}).Where("team_id = ?", team.ID).Count(&coachCount)
		if coachCount == 0 {
			coach, err := s.p.Coach(ctx, extID)
			if err != nil {
				syncErrors = append(syncErrors, fmt.Errorf("sync coach team %s: %w", team.Name, err))
			} else if coach != nil {
				cExtID := coach.ExternalID
				cRecord := models.Coach{
					TeamID:        team.ID,
					Name:          coach.Name,
					Nationality:   coach.Nationality,
					PhotoURL:      coach.PhotoURL,
					ExternalAPIID: &cExtID,
				}
				var existing []models.Coach
				s.db.WithContext(ctx).Where("external_api_id = ?", cExtID).Limit(1).Find(&existing)
				if len(existing) == 0 {
					if err := s.db.WithContext(ctx).Create(&cRecord).Error; err == nil {
						result.CoachesCount++
					}
				} else {
					s.db.WithContext(ctx).Model(&existing[0]).Updates(map[string]any{
						"name":        cRecord.Name,
						"nationality": cRecord.Nationality,
						"photo_url":   cRecord.PhotoURL,
					})
				}
			}
		}

		// 2. Sync Squad (ถ้ามีรายชื่ออยู่แล้วให้ข้ามเพื่อประหยัดโควตา)
		var playerCount int64
		s.db.WithContext(ctx).Model(&models.Player{}).Where("team_id = ?", team.ID).Count(&playerCount)
		if playerCount == 0 {
			players, err := s.p.Squad(ctx, extID)
			if err != nil {
				syncErrors = append(syncErrors, fmt.Errorf("sync squad team %s: %w", team.Name, err))
			} else {
				for _, p := range players {
					pExtID := p.ExternalID
					pRecord := models.Player{
						TeamID:        team.ID,
						Name:          p.Name,
						Number:        p.Number,
						Position:      p.Position,
						Age:           p.Age,
						PhotoURL:      p.PhotoURL,
						ExternalAPIID: &pExtID,
					}
					var existing []models.Player
					s.db.WithContext(ctx).Where("external_api_id = ?", pExtID).Limit(1).Find(&existing)
					if len(existing) == 0 {
						if err := s.db.WithContext(ctx).Create(&pRecord).Error; err == nil {
							result.PlayersCount++
						}
					} else {
						s.db.WithContext(ctx).Model(&existing[0]).Updates(map[string]any{
							"name":      pRecord.Name,
							"number":    pRecord.Number,
							"position":  pRecord.Position,
							"age":       pRecord.Age,
							"photo_url": pRecord.PhotoURL,
						})
					}
				}
				s.log.Info("team squad synced", "team", team.Name, "players", len(players))
			}
		} else {
			s.log.Info("team squad already present, skipped", "team", team.Name, "players", playerCount)
		}
	}

	return result, errors.Join(syncErrors...)
}

