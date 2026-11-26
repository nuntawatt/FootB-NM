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

type leagueInput struct {
	Name          string  `json:"name"`
	ShortName     string  `json:"short_name"`
	Country       string  `json:"country"`
	LogoURL       string  `json:"logo_url"`
	ExternalAPIID *string `json:"external_api_id"`
}

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

func (h *LeagueHandler) Create(c fiber.Ctx) error {
	var in leagueInput
	if err := c.Bind().Body(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	if in.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}
	league := models.League{
		Name: in.Name, ShortName: in.ShortName, Country: in.Country,
		LogoURL: in.LogoURL, ExternalAPIID: in.ExternalAPIID,
	}
	if err := h.db.Create(&league).Error; err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(league)
}

func (h *LeagueHandler) Update(c fiber.Ctx) error {
	league, err := h.find(c)
	if err != nil {
		return err
	}
	var in leagueInput
	if err := c.Bind().Body(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	if in.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}
	league.Name, league.ShortName, league.Country = in.Name, in.ShortName, in.Country
	league.LogoURL, league.ExternalAPIID = in.LogoURL, in.ExternalAPIID
	if err := h.db.Save(league).Error; err != nil {
		return err
	}
	return c.JSON(league)
}

func (h *LeagueHandler) Delete(c fiber.Ctx) error {
	league, err := h.find(c)
	if err != nil {
		return err
	}
	if err := h.db.Delete(league).Error; err != nil { // soft delete
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
