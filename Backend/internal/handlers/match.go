package handlers

import (
	"errors"
	"strconv"
	"time"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MatchHandler struct{ db *gorm.DB }

func NewMatchHandler(db *gorm.DB) *MatchHandler { return &MatchHandler{db: db} }

func (h *MatchHandler) List(c fiber.Ctx) error {
	var matches []models.Match
	q := h.db.Preload("HomeTeam").Preload("AwayTeam").Preload("League")

	if leagueID := c.Query("league_id"); leagueID != "" {
		if lid, err := uuid.Parse(leagueID); err == nil {
			q = q.Where("league_id = ?", lid)
		}
	}

	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	if dateStr := c.Query("date"); dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			end := start.Add(24 * time.Hour)
			q = q.Where("match_date >= ? AND match_date < ?", start, end)
		}
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	q = q.Limit(limit)

	if pageS := c.Query("page"); pageS != "" {
		if p, err := strconv.Atoi(pageS); err == nil && p > 1 {
			q = q.Offset((p - 1) * limit)
		}
	}

	if err := q.Order("match_date desc").Find(&matches).Error; err != nil {
		return err
	}
	return c.JSON(matches)
}

func (h *MatchHandler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	var match models.Match
	if err := h.db.Preload("HomeTeam").Preload("AwayTeam").Preload("League").First(&match, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "match not found")
		}
		return err
	}
	return c.JSON(match)
}
