package handlers

import (
	"strings"
	"time"

	"football/internal/config"
	"football/internal/middleware"
	"football/internal/models"
	"football/internal/validator"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{db: db, cfg: cfg}
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func toUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

func (h *UserHandler) generateToken(user models.User) (string, error) {
	claims := middleware.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)), 
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

// Register สมัครสมาชิกใหม่
func (h *UserHandler) Register(c fiber.Ctx) error {
	var in validator.RegisterInput
	if err := c.Bind().Body(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// ใช้ validator ตรวจสอบข้อมูลผู้ใช้อย่างละเอียด
	if err := in.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var existing []models.User
	if err := h.db.Where("email = ? OR username = ?", in.Email, in.Username).Limit(1).Find(&existing).Error; err == nil && len(existing) > 0 {
		if strings.EqualFold(existing[0].Email, in.Email) {
			return fiber.NewError(fiber.StatusConflict, "email is already registered")
		}
		return fiber.NewError(fiber.StatusConflict, "username is already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to hash password")
	}

	avatar := in.AvatarURL
	if avatar == "" {
		avatar = "https://i.pinimg.com/1200x/81/8b/7e/818b7e92210d0afdbbb2f1f84247e1c4.jpg"
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
		"user":   toUserResponse(user),
	})
}

// Login เข้าสู่ระบบ
func (h *UserHandler) Login(c fiber.Ctx) error {
	var in validator.LoginInput
	if err := c.Bind().Body(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := in.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	identifier := in.Email
	var users []models.User
	err := h.db.Where("lower(email) = ? OR username = ?", strings.ToLower(identifier), identifier).Limit(1).Find(&users).Error
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid email/username or password")
	}
	user := users[0]

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
		"user":   toUserResponse(user),
	})
}

// Profile คืนข้อมูลโปรไฟล์ของผู้ใช้ที่ล็อกอินอยู่
func (h *UserHandler) Profile(c fiber.Ctx) error {
	val := c.Locals("userID")
	if val == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	userID := val.(uuid.UUID)

	var users []models.User
	if err := h.db.Limit(1).Find(&users, "id = ?", userID).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(toUserResponse(users[0]))
}

// UpdateProfile แก้ไขข้อมูลโปรไฟล์
func (h *UserHandler) UpdateProfile(c fiber.Ctx) error {
	val := c.Locals("userID")
	if val == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	userID := val.(uuid.UUID)

	var in validator.UpdateProfileInput
	if err := c.Bind().Body(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := in.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var users []models.User
	if err := h.db.Limit(1).Find(&users, "id = ?", userID).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	user := users[0]

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
		"user":    toUserResponse(user),
	})
}
