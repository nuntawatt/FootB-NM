package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Base
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	AvatarURL    string `json:"avatar_url"`
	Role         string `gorm:"type:varchar(10);default:user;not null" json:"role"`
}

type UserFavoriteTeam struct {
	UserID    uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"user_id"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"team_id"`
	CreatedAt time.Time `json:"created_at"`
}
