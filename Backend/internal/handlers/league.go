package handlers

import (
	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeagueHandler struct{ db *gorm.DB }

func NewLeagueHandler(db *gorm.DB) *LeagueHandler { return &LeagueHandler{db: db} }

// List คืนรายการลีกทั้งหมด
func (h *LeagueHandler) List(c fiber.Ctx) error {
	var leagues []models.League
	if err := h.db.Order("name").Find(&leagues).Error; err != nil {
		return err
	}
	return c.JSON(leagues)
}

// Get คืนข้อมูลลีกตาม id
func (h *LeagueHandler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var leagues []models.League
	if err := h.db.Limit(1).Find(&leagues, "id = ?", id).Error; err != nil {
		return err
	}
	if len(leagues) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "league not found")
	}
	return c.JSON(leagues[0])
}
