package models

import "github.com/google/uuid"

type Coach struct {
	Base
	TeamID        uuid.UUID `gorm:"type:uuid;not null;index" json:"team_id"`
	Team          *Team     `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Name          string    `gorm:"not null" json:"name"`
	Nationality   string    `json:"nationality"`
	PhotoURL      string    `json:"photo_url"`
	ExternalAPIID *string   `gorm:"uniqueIndex" json:"-"`
}
