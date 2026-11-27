package validator

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,10}$`)
)

type RegisterInput struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	AvatarURL       string `json:"avatar_url"`
}

func (r *RegisterInput) Sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.AvatarURL = strings.TrimSpace(r.AvatarURL)
}

func (r *RegisterInput) Validate() error {
	r.Sanitize()

	if r.Username == "" {
		return errors.New("username is required")
	}
	if len(r.Username) < 3 || len(r.Username) > 30 {
		return errors.New("username must be between 3 and 10 characters")
	}
	if !usernameRegex.MatchString(r.Username) {
		return errors.New("username can only contain letters, numbers, and underscores (no spaces or special characters)")
	}

	if r.Email == "" {
		return errors.New("email is required")
	}
	if len(r.Email) > 255 {
		return errors.New("email cannot exceed 255 characters")
	}
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email address format (e.g. user@example.com)")
	}

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

	if r.ConfirmPassword == "" {
		return errors.New("confirm_password is required")
	}
	if r.Password != r.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	if r.AvatarURL != "" {
		u, err := url.ParseRequestURI(r.AvatarURL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return errors.New("avatar_url must be a valid HTTP or HTTPS URL")
		}
	}

	return nil
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginInput) Sanitize() {
	l.Email = strings.TrimSpace(l.Email)
}

func (l *LoginInput) Validate() error {
	l.Sanitize()
	if l.Email == "" {
		return errors.New("email or username is required")
	}
	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type UpdateProfileInput struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

func (u *UpdateProfileInput) Sanitize() {
	u.Username = strings.TrimSpace(u.Username)
	u.AvatarURL = strings.TrimSpace(u.AvatarURL)
}

func (u *UpdateProfileInput) Validate() error {
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
