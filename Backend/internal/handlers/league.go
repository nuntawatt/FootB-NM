package handlers

import (
	"errors"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeagueHandler struct{ db *gorm.DB }

func NewLeagueHandler(db *gorm.DB) *LeagueHandler { return &LeagueHandler{db: db} }

// หา league จาก :id ใช้ร่วมกันหลาย endpoint
func (h *LeagueHandler) find(c fiber.Ctx) (*models.League, error) {
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

func (h *LeagueHandler) List(c fiber.Ctx) error {
	var leagues []models.League
	if err := h.db.Order("name").Find(&leagues).Error; err != nil {
		return err
	}
	return c.JSON(leagues)
}

func (h *LeagueHandler) Get(c fiber.Ctx) error {
	league, err := h.find(c)
	if err != nil {
		return err
	}
	return c.JSON(league)
}
