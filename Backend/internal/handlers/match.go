package handlers

import (
	"strconv"
	"strings"
	"time"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var availableSeasons = []int{2022, 2023, 2024, 2025, 2026}

type MatchHandler struct{ db *gorm.DB }

func NewMatchHandler(db *gorm.DB) *MatchHandler { return &MatchHandler{db: db} }

// List คืนรายการแมตช์ตาม filter (league_id, status, season, year, team_id, date, page, limit, group_by)
func (h *MatchHandler) List(c fiber.Ctx) error {
	q := h.db.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("League")

	q = applyMatchFilters(q, c)

	// หากมีการกรองเฉพาะปีหรือทีม ให้ดึงได้มากขึ้น (default 300, max 1000)
	defaultLimit := 50
	hasSpecificFilter := c.Query("season") != "" || c.Query("year") != "" || c.Query("group_by") != "" || c.Query("team_id") != ""
	if hasSpecificFilter {
		defaultLimit = 300
	}
	limit := parseLimitParam(c.Query("limit"), defaultLimit, 1000)
	q = q.Limit(limit)

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 1 {
		q = q.Offset((page - 1) * limit)
	}

	var matches []models.Match
	if err := q.Order("match_date desc").Find(&matches).Error; err != nil {
		return err
	}

	// รองรับการจัดกลุ่มแมตช์ตามปี (?group_by=year) หรือตามฤดูกาล (?group_by=season)
	groupBy := strings.ToLower(c.Query("group_by"))
	if groupBy == "year" {
		grouped := make(map[string][]models.Match)
		for _, m := range matches {
			key := strconv.Itoa(m.Year)
			grouped[key] = append(grouped[key], m)
		}
		return c.JSON(grouped)
	} else if groupBy == "season" {
		grouped := make(map[string][]models.Match)
		for _, m := range matches {
			key := m.Season
			grouped[key] = append(grouped[key], m)
		}
		return c.JSON(grouped)
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
	// กรองตามลีก League ID
	if raw := c.Query("league_id"); raw != "" {
		if lid, err := uuid.Parse(raw); err == nil {
			q = q.Where("league_id = ?", lid)
		}
	}

	// กรองตามสถานะการแข่งขัน status: finished, scheduled, live, postponed
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	// กรองแมตช์ของทีมใดทีมหนึ่ง ไม่ว่าจะเป็นทีมเหย้าหรือเยือน
	if raw := c.Query("team_id"); raw != "" {
		if tid, err := uuid.Parse(raw); err == nil {
			q = q.Where("home_team_id = ? OR away_team_id = ?", tid, tid)
		}
	}

	// กรองคู่แข่งขันเฉพาะเจาะจง home_team_id / away_team_id
	if raw := c.Query("home_team_id"); raw != "" {
		if hid, err := uuid.Parse(raw); err == nil {
			q = q.Where("home_team_id = ?", hid)
		}
	}
	if raw := c.Query("away_team_id"); raw != "" {
		if aid, err := uuid.Parse(raw); err == nil {
			q = q.Where("away_team_id = ?", aid)
		}
	}

	// กรองตามปีปฏิทิน เช่น ?year=2022, 2023, 2024, 2025, 2026 -> 1 ม.ค. ถึง 31 ธ.ค.
	if raw := c.Query("year"); raw != "" {
		if year, err := strconv.Atoi(raw); err == nil {
			start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
			q = q.Where("match_date >= ? AND match_date < ?", start, end)
		}
	}

	// กรองตามฤดูกาลฟุตบอล เช่น ?season=2024 หรือ ?season=2024/2025 -> 1 ก.ค. 2024 ถึง 1 ก.ค. 2025
	if raw := c.Query("season"); raw != "" {
		seasonYearStr := strings.Split(raw, "/")[0]
		if year, err := strconv.Atoi(seasonYearStr); err == nil {
			start, end := seasonDateRange(year)
			q = q.Where("match_date >= ? AND match_date < ?", start, end)
		}
	}

	// กรองตามวันที่ระบุ YYYY-MM-DD
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
