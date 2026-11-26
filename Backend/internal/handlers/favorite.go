package handlers

import (
	"time"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FavoriteHandler struct{ db *gorm.DB }

func NewFavoriteHandler(db *gorm.DB) *FavoriteHandler {
	return &FavoriteHandler{db: db}
}

func (h *FavoriteHandler) getEffectiveUserID(c fiber.Ctx) uuid.UUID {
	if val := c.Locals("userID"); val != nil {
		if uid, ok := val.(uuid.UUID); ok {
			return uid
		}
	}
	if headerVal := c.Get("X-User-ID"); headerVal != "" {
		if uid, err := uuid.Parse(headerVal); err == nil {
			return uid
		}
	}
	defaultID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	return defaultID
}

// List คืนรายการทีมโปรด
func (h *FavoriteHandler) List(c fiber.Ctx) error {
	userID := h.getEffectiveUserID(c)

	var teamIDs []uuid.UUID
	err := h.db.Model(&models.UserFavoriteTeam{}).
		Where("user_id = ?", userID).
		Pluck("team_id", &teamIDs).Error
	if err != nil {
		return err
	}

	if len(teamIDs) == 0 {
		return c.JSON([]models.Team{})
	}

	var teams []models.Team
	if err := h.db.Where("id IN ?", teamIDs).Order("name").Find(&teams).Error; err != nil {
		return err
	}
	return c.JSON(teams)
}

// Add เพิ่มทีมลงในทีมโปรด
func (h *FavoriteHandler) Add(c fiber.Ctx) error {
	userID := h.getEffectiveUserID(c)
	teamID, err := uuid.Parse(c.Params("team_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid team_id")
	}

	var count int64
	h.db.Model(&models.Team{}).Where("id = ?", teamID).Count(&count)
	if count == 0 {
		return fiber.NewError(fiber.StatusNotFound, "team not found")
	}

	fav := models.UserFavoriteTeam{
		UserID:    userID,
		TeamID:    teamID,
		CreatedAt: time.Now(),
	}

	if err := h.db.Where("user_id = ? AND team_id = ?", userID, teamID).FirstOrCreate(&fav).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "team added to favorites",
		"team_id": teamID,
	})
}

// Remove ลบทีมออกจากทีมโปรด
func (h *FavoriteHandler) Remove(c fiber.Ctx) error {
	userID := h.getEffectiveUserID(c)
	teamID, err := uuid.Parse(c.Params("team_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid team_id")
	}

	res := h.db.Where("user_id = ? AND team_id = ?", userID, teamID).Delete(&models.UserFavoriteTeam{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "favorite not found")
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "team removed from favorites",
		"team_id": teamID,
	})
}
