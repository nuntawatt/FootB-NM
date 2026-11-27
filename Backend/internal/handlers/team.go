package handlers

import (
	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamHandler struct{ db *gorm.DB }

func NewTeamHandler(db *gorm.DB) *TeamHandler { return &TeamHandler{db: db} }

// List คืนรายการทีม สามารถ filter ?league_id= ได้
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

// Get คืนข้อมูลทีมตาม id
func (h *TeamHandler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var teams []models.Team
	if err := h.db.Limit(1).Find(&teams, "id = ?", id).Error; err != nil {
		return err
	}
	if len(teams) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "team not found")
	}
	return c.JSON(teams[0])
}
