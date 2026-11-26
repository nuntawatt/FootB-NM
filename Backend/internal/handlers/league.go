// Package handlers contains HTTP handler types for the football API.
// Each handler type wraps a *gorm.DB and exposes methods that can be
// registered directly as Fiber route handlers.
package handlers

import (
	"errors"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LeagueHandler handles HTTP requests for the /leagues resource.
type LeagueHandler struct{ db *gorm.DB }

// NewLeagueHandler creates a LeagueHandler backed by the given database.
func NewLeagueHandler(db *gorm.DB) *LeagueHandler { return &LeagueHandler{db: db} }

// List handles GET /leagues.
// Returns all leagues ordered by name.
func (h *LeagueHandler) List(c fiber.Ctx) error {
	var leagues []models.League
	if err := h.db.Order("name").Find(&leagues).Error; err != nil {
		return err
	}
	return c.JSON(leagues)
}

// Get handles GET /leagues/:id.
// Returns a single league by its UUID, or 404 if not found.
func (h *LeagueHandler) Get(c fiber.Ctx) error {
	league, err := h.findByID(c)
	if err != nil {
		return err
	}
	return c.JSON(league)
}

// findByID is a shared helper that parses the :id param and loads the league.
func (h *LeagueHandler) findByID(c fiber.Ctx) (*models.League, error) {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var league models.League
	if err := h.db.First(&league, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "league not found")
		}
		return nil, err
	}
	return &league, nil
}
