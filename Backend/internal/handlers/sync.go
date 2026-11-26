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

// SyncHandler handles on-demand data synchronisation from API-Football.
type SyncHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewSyncHandler creates a SyncHandler with the given database and config.
func NewSyncHandler(db *gorm.DB, cfg *config.Config) *SyncHandler {
	return &SyncHandler{db: db, cfg: cfg}
}

// syncRequest is the JSON body accepted by SyncFixtures.
type syncRequest struct {
	Season int    `json:"season"`           // season year, e.g. 2024
	From   string `json:"from,omitempty"`   // optional start date YYYY-MM-DD
	To     string `json:"to,omitempty"`     // optional end date   YYYY-MM-DD
}

// SyncFixtures handles POST /sync/fixtures.
//
// It fetches fixtures from API-Football for all leagues that have an
// external_api_id set, then upserts the results into the database.
//
// Request body (JSON):
//
//	{ "season": 2024, "from": "2024-08-01", "to": "2024-12-31" }
//
// "from" and "to" are optional. When omitted the full season is synced.
// "season" falls back to FOOTBALL_SEASON from the environment when omitted.
func (h *SyncHandler) SyncFixtures(c fiber.Ctx) error {
	var req syncRequest
	_ = c.Bind().Body(&req)

	// Fall back to the configured default season.
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

// parseDateRange parses optional YYYY-MM-DD strings into time.Time values.
// Dates are interpreted in Asia/Bangkok time to match API-Football's timezone.
// Returns zero values when the strings are empty or malformed.
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
