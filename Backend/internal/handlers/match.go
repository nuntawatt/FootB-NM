package handlers

import (
	"strconv"
	"time"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var availableSeasons = []int{2022, 2023, 2024, 2025, 2026}

type MatchHandler struct{ db *gorm.DB }

func NewMatchHandler(db *gorm.DB) *MatchHandler { return &MatchHandler{db: db} }

// List คืนรายการแมตช์ตาม filter (league_id, status, season, year, date, page, limit)
func (h *MatchHandler) List(c fiber.Ctx) error {
	q := h.db.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("League")

	q = applyMatchFilters(q, c)

	limit := parseLimitParam(c.Query("limit"), 50, 200)
	q = q.Limit(limit)

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 1 {
		q = q.Offset((page - 1) * limit)
	}

	var matches []models.Match
	if err := q.Order("match_date desc").Find(&matches).Error; err != nil {
		return err
	}
	return c.JSON(matches)
}

// Get คืนข้อมูลแมตช์เดี่ยวตาม id
func (h *MatchHandler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var matches []models.Match
	err = h.db.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("League").
		Limit(1).
		Find(&matches, "id = ?", id).Error

	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "match not found")
	}
	return c.JSON(matches[0])
}

// Seasons คืนรายการปีฤดูกาลและจำนวนแมตช์ในระบบ
func (h *MatchHandler) Seasons(c fiber.Ctx) error {
	type SeasonItem struct {
		Season       int    `json:"season"`
		Label        string `json:"label"`
		MatchesCount int64  `json:"matches_count"`
	}

	results := make([]SeasonItem, 0, len(availableSeasons))

	for _, year := range availableSeasons {
		start, end := seasonDateRange(year)

		var count int64
		h.db.Model(&models.Match{}).
			Where("match_date >= ? AND match_date < ?", start, end).
			Count(&count)

		results = append(results, SeasonItem{
			Season:       year,
			Label:        strconv.Itoa(year) + "/" + strconv.Itoa(year+1),
			MatchesCount: count,
		})
	}

	return c.JSON(results)
}

func applyMatchFilters(q *gorm.DB, c fiber.Ctx) *gorm.DB {
	if raw := c.Query("league_id"); raw != "" {
		if lid, err := uuid.Parse(raw); err == nil {
			q = q.Where("league_id = ?", lid)
		}
	}

	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	// กรองตามฤดูกาล (เช่น 2024 -> 1 ก.ค. 2024 ถึง 1 ก.ค. 2025)
	if raw := c.Query("season"); raw != "" {
		if year, err := strconv.Atoi(raw); err == nil {
			start, end := seasonDateRange(year)
			q = q.Where("match_date >= ? AND match_date < ?", start, end)
		}
	}

	// กรองตามปีปฏิทิน (1 ม.ค. ถึง 31 ธ.ค.)
	if raw := c.Query("year"); raw != "" {
		if year, err := strconv.Atoi(raw); err == nil {
			start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
			q = q.Where("match_date >= ? AND match_date < ?", start, end)
		}
	}

	// กรองตามวันที่ระบุ (YYYY-MM-DD)
	if raw := c.Query("date"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			end := start.Add(24 * time.Hour)
			q = q.Where("match_date >= ? AND match_date < ?", start, end)
		}
	}

	return q
}

func seasonDateRange(year int) (start, end time.Time) {
	start = time.Date(year, time.July, 1, 0, 0, 0, 0, time.UTC)
	end = time.Date(year+1, time.July, 1, 0, 0, 0, 0, time.UTC)
	return
}

func parseLimitParam(raw string, defaultVal, max int) int {
	if raw == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return defaultVal
	}
	if n > max {
		return max
	}
	return n
}
