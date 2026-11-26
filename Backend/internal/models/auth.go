package models

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
)

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
}

func (r *RegisterRequest) Sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.AvatarURL = strings.TrimSpace(r.AvatarURL)
}

// Validate ตรวจสอบความถูกต้องของข้อมูลสมัครสมาชิกอย่างละเอียด
func (r *RegisterRequest) Validate() error {
	r.Sanitize()

	// 1. ตรวจสอบ Username
	if r.Username == "" {
		return errors.New("username is required")
	}
	if len(r.Username) < 3 || len(r.Username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}
	if !usernameRegex.MatchString(r.Username) {
		return errors.New("username can only contain letters, numbers, and underscores (no spaces or special characters)")
	}

	// 2. ตรวจสอบ Email
	if r.Email == "" {
		return errors.New("email is required")
	}
	if len(r.Email) > 255 {
		return errors.New("email cannot exceed 255 characters")
	}
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email address format (e.g. user@example.com)")
	}

	// 3. ตรวจสอบ Password
	if r.Password == "" {
		return errors.New("password is required")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if len(r.Password) > 72 {
		return errors.New("password cannot exceed 72 characters")
	}

	var hasLetter, hasNumber bool
	for _, ch := range r.Password {
		switch {
		case unicode.IsLetter(ch):
			hasLetter = true
		case unicode.IsNumber(ch):
			hasNumber = true
		}
	}
	if !hasLetter || !hasNumber {
		return errors.New("password must contain at least one letter and one number")
	}

	// 4. ตรวจสอบ AvatarURL (ถ้ามี)
	if r.AvatarURL != "" {
		u, err := url.ParseRequestURI(r.AvatarURL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return errors.New("avatar_url must be a valid HTTP or HTTPS URL")
		}
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"` // สามารถเป็นอีเมลหรือชื่อผู้ใช้ได้
	Password string `json:"password"`
}

func (l *LoginRequest) Sanitize() {
	l.Email = strings.TrimSpace(l.Email)
}

func (l *LoginRequest) Validate() error {
	l.Sanitize()
	if l.Email == "" {
		return errors.New("email or username is required")
	}
	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type UpdateProfileRequest struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

func (u *UpdateProfileRequest) Sanitize() {
	u.Username = strings.TrimSpace(u.Username)
	u.AvatarURL = strings.TrimSpace(u.AvatarURL)
}

func (u *UpdateProfileRequest) Validate() error {
	u.Sanitize()
	if u.Username == "" && u.AvatarURL == "" {
		return errors.New("at least one field (username or avatar_url) must be provided to update")
	}
	if u.Username != "" {
		if len(u.Username) < 3 || len(u.Username) > 30 {
			return errors.New("username must be between 3 and 30 characters")
		}
		if !usernameRegex.MatchString(u.Username) {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}
	if u.AvatarURL != "" {
		parsed, err := url.ParseRequestURI(u.AvatarURL)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return errors.New("avatar_url must be a valid HTTP or HTTPS URL")
		}
	}
	return nil
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserResponse(u User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
