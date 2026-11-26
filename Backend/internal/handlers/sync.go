package handlers

import (
	"log/slog"
	"time"

	"football/internal/config"
	"football/internal/ingest"
	"football/internal/provider/apifootball"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type SyncHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewSyncHandler(db *gorm.DB, cfg *config.Config) *SyncHandler {
	return &SyncHandler{db: db, cfg: cfg}
}

type syncRequest struct {
	Season int    `json:"season"`
	From   string `json:"from"`
	To     string `json:"to"`
}

func (h *SyncHandler) SyncFixtures(c fiber.Ctx) error {
	var in syncRequest
	_ = c.Bind().Body(&in)

	if in.Season == 0 {
		in.Season = h.cfg.Season
	}
	if in.Season == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "season is required (e.g. 2024)")
	}

	var from, to time.Time
	loc, _ := time.LoadLocation(apifootball.DefaultTimezone)
	if in.From != "" {
		if t, err := time.ParseInLocation("2006-01-02", in.From, loc); err == nil {
			from = t
		}
	}
	if in.To != "" {
		if t, err := time.ParseInLocation("2006-01-02", in.To, loc); err == nil {
			to = t
		}
	}

	client := apifootball.New(h.cfg.APIFootballKey, apifootball.WithBaseURL(h.cfg.APIFootballBaseURL))
	svc := ingest.New(h.db, client, slog.Default())

	res, err := svc.SyncFixtures(c.Context(), in.Season, from, to)
	if err != nil && res.MatchesCreated == 0 && res.MatchesUpdated == 0 {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	dayQuota, minQuota := client.Quota()
	return c.JSON(fiber.Map{
		"status":          "success",
		"season":          in.Season,
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
