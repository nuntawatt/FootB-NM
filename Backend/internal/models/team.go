package models

import "github.com/google/uuid"

type Team struct {
	Base
	LeagueID      uuid.UUID `gorm:"type:uuid;not null;index" json:"league_id"`
	League        *League   `gorm:"foreignKey:LeagueID" json:"league,omitempty"`
	Name          string    `gorm:"not null" json:"name"`
	ShortName     string    `json:"short_name"`
	LogoURL       string    `json:"logo_url"`
	ExternalAPIID *string   `gorm:"uniqueIndex" json:"-"` // id สำหรับซิงก์ external API (ไม่แสดงใน JSON)
}
