package handlers

import (
	"errors"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamHandler struct{ db *gorm.DB }

func NewTeamHandler(db *gorm.DB) *TeamHandler { return &TeamHandler{db: db} }

func (h *TeamHandler) List(c fiber.Ctx) error {
	var teams []models.Team
	q := h.db.Model(&models.Team{})
	if leagueID := c.Query("league_id"); leagueID != "" {
		if lid, err := uuid.Parse(leagueID); err == nil {
			q = q.Where("league_id = ?", lid)
		}
	}
	if err := q.Order("name").Find(&teams).Error; err != nil {
		return err
	}
	return c.JSON(teams)
}

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
