package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// Protected ตรวจสอบ JWT Bearer Token หากไม่ถูกต้องจะคืน 401 Unauthorized
func Protected(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization format; expected 'Bearer <token>'")
		}

		tokenString := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// OptionalAuth อ่าน Token หากส่งมา แต่ถ้าไม่ส่งก็ไม่ block
func OptionalAuth(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				claims := &JWTClaims{}
				token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (any, error) {
					return []byte(jwtSecret), nil
				})
				if err == nil && token.Valid {
					c.Locals("userID", claims.UserID)
					c.Locals("username", claims.Username)
					c.Locals("userRole", claims.Role)
				}
			}
		}
		return c.Next()
	}
}
