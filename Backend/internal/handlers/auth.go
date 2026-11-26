package handlers

import (
	"errors"
	"strings"
	"time"

	"football/internal/config"
	"football/internal/middleware"
	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

func (h *AuthHandler) generateToken(user models.User) (string, error) {
	claims := middleware.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // อายุ 7 วัน
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

// Register สมัครสมาชิกใหม่พร้อม Validation ละเอียด
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var in models.RegisterRequest
	if err := c.Bind().Body(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Validate ข้อมูลนำเข้าอย่างละเอียด
	if err := in.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// ตรวจสอบชื่อผู้ใช้หรืออีเมลซ้ำในฐานข้อมูล
	var existing []models.User
	if err := h.db.Where("email = ? OR username = ?", in.Email, in.Username).Limit(1).Find(&existing).Error; err == nil && len(existing) > 0 {
		if strings.EqualFold(existing[0].Email, in.Email) {
			return fiber.NewError(fiber.StatusConflict, "email is already registered")
		}
		return fiber.NewError(fiber.StatusConflict, "username is already taken")
	}

	// แฮชรหัสผ่านด้วย bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to hash password")
	}

	avatar := in.AvatarURL
	if avatar == "" {
		avatar = "https://ui-avatars.com/api/?name=" + in.Username + "&background=0F172A&color=38BDF8&size=256"
	}

	user := models.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: string(hash),
		AvatarURL:    avatar,
		Role:         "user",
	}

	if err := h.db.Create(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create user account")
	}

	token, err := h.generateToken(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate authentication token")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"token":  token,
		"user":   models.ToUserResponse(user),
	})
}

// Login เข้าสู่ระบบพร้อม Validation
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var in models.LoginRequest
	if err := c.Bind().Body(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := in.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	identifier := in.Email
	var user models.User
	err := h.db.Where("lower(email) = ? OR username = ?", strings.ToLower(identifier), identifier).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid email/username or password")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid email/username or password")
	}

	token, err := h.generateToken(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate authentication token")
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"token":  token,
		"user":   models.ToUserResponse(user),
	})
}

// Profile คืนข้อมูลโปรไฟล์ของผู้ใช้ที่ล็อกอินอยู่
func (h *AuthHandler) Profile(c fiber.Ctx) error {
	val := c.Locals("userID")
	if val == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	userID := val.(uuid.UUID)

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(models.ToUserResponse(user))
}

// UpdateProfile แก้ไขข้อมูลโปรไฟล์พร้อม Validation
func (h *AuthHandler) UpdateProfile(c fiber.Ctx) error {
	val := c.Locals("userID")
	if val == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	userID := val.(uuid.UUID)

	var in models.UpdateProfileRequest
	if err := c.Bind().Body(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := in.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	updates := make(map[string]any)
	if in.Username != "" && in.Username != user.Username {
		var existing []models.User
		if err := h.db.Where("username = ? AND id != ?", in.Username, userID).Limit(1).Find(&existing).Error; err == nil && len(existing) > 0 {
			return fiber.NewError(fiber.StatusConflict, "username is already taken")
		}
		updates["username"] = in.Username
	}
	if in.AvatarURL != "" {
		updates["avatar_url"] = in.AvatarURL
	}

	if len(updates) > 0 {
		if err := h.db.Model(&user).Updates(updates).Error; err != nil {
			return err
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "profile updated successfully",
		"user":    models.ToUserResponse(user),
	})
}
