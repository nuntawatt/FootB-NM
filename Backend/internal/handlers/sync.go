package handlers

import (
	"log/slog"
	"time"

	"football/internal/redis"
	"football/internal/config"
	"football/internal/ingest"
	"football/internal/models"
	"football/internal/provider/apifootball"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type SyncHandler struct {
	db  *gorm.DB
	cfg *config.Config
	cache *cache.Cache
}

func NewSyncHandler(db *gorm.DB, cfg *config.Config, cache *cache.Cache) *SyncHandler {
	return &SyncHandler{db: db, cfg: cfg, cache: cache}
}

// ดึงข้อมูลจาก API-Football ลงฐานข้อมูลตามฤดูกาลที่ระบุ
func (h *SyncHandler) SyncFixtures(c fiber.Ctx) error {
	var req models.SyncRequest
	_ = c.Bind().Body(&req)

	if req.Season == 0 {
		req.Season = h.cfg.Season
	}
	if req.Season == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "season is required (e.g. 2024)")
	}

	from, to := parseDateRange(req.From, req.To)

	client := apifootball.New(
		h.cfg.APIFootballKey,
		apifootball.WithBaseURL(h.cfg.APIFootballBaseURL),
	)
	svc := ingest.New(h.db, client, slog.Default())

	res, err := svc.SyncFixtures(c.Context(), req.Season, from, to)
	if err != nil && res.MatchesCreated == 0 && res.MatchesUpdated == 0 {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// เมื่อ Sync สำเร็จ ให้ลบ Cache ตารางคะแนน เพื่อให้คำนวณใหม่
	h.cache.DeletePrefix(c.Context(), "football:standings:")

	// อ่าน Quota ที่เหลือของ API-Football
	dayQuota, minQuota := client.Quota()

	return c.JSON(fiber.Map{
		"status":          "success",
		"season":          req.Season,
		"leagues_synced":  res.Leagues,
		"teams_created":   res.TeamsCreated,
		"matches_created": res.MatchesCreated,
		"matches_updated": res.MatchesUpdated,
		"matches_skipped": res.Skipped,
		"quota_remaining": fiber.Map{
			"day":    dayQuota,
			"minute": minQuota,
		},
	})
}

// SyncSquads ดึงโค้ชและนักเตะของทุกทีมลง Database
func (h *SyncHandler) SyncSquads(c fiber.Ctx) error {
	client := apifootball.New(
		h.cfg.APIFootballKey,
		apifootball.WithBaseURL(h.cfg.APIFootballBaseURL),
	)
	svc := ingest.New(h.db, client, slog.Default())

	res, err := svc.SyncSquadsAndCoaches(c.Context())
	if err != nil && res.CoachesCount == 0 && res.PlayersCount == 0 {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	dayQuota, minQuota := client.Quota()
	return c.JSON(fiber.Map{
		"status":          "success",
		"teams_checked":   res.TeamsCount,
		"coaches_synced":  res.CoachesCount,
		"players_synced":  res.PlayersCount,
		"quota_remaining": fiber.Map{
			"day":    dayQuota,
			"minute": minQuota,
		},
	})
}

func parseDateRange(fromStr, toStr string) (from, to time.Time) {
	loc, _ := time.LoadLocation(apifootball.DefaultTimezone)

	if fromStr != "" {
		if t, err := time.ParseInLocation("2006-01-02", fromStr, loc); err == nil {
			from = t
		}
	}
	if toStr != "" {
		if t, err := time.ParseInLocation("2006-01-02", toStr, loc); err == nil {
			to = t
		}
	}
	return
}
