package handlers

import (
	"errors"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TeamHandler handles HTTP requests for the /teams resource.
type TeamHandler struct{ db *gorm.DB }

// NewTeamHandler creates a TeamHandler backed by the given database.
func NewTeamHandler(db *gorm.DB) *TeamHandler { return &TeamHandler{db: db} }

// List handles GET /teams.
// Accepts an optional ?league_id=<uuid> query parameter to filter by league.
// Results are ordered alphabetically by name.
func (h *TeamHandler) List(c fiber.Ctx) error {
	q := h.db.Model(&models.Team{})

	if raw := c.Query("league_id"); raw != "" {
		if lid, err := uuid.Parse(raw); err == nil {
			q = q.Where("league_id = ?", lid)
		}
	}

	var teams []models.Team
	if err := q.Order("name").Find(&teams).Error; err != nil {
		return err
	}
	return c.JSON(teams)
}

// Get handles GET /teams/:id.
// Returns a single team by its UUID, or 404 if not found.
func (h *TeamHandler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var team models.Team
	if err := h.db.First(&team, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "team not found")
		}
		return err
	}
	return c.JSON(team)
}
