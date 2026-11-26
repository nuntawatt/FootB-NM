package models

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	Base
	LeagueID         uuid.UUID `gorm:"type:uuid;not null;index" json:"league_id"`
	League           *League   `gorm:"foreignKey:LeagueID" json:"league,omitempty"`
	HomeTeamID       uuid.UUID `gorm:"type:uuid;not null;index" json:"home_team_id"`
	HomeTeam         *Team     `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	AwayTeamID       uuid.UUID `gorm:"type:uuid;not null;index" json:"away_team_id"`
	AwayTeam         *Team     `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	MatchDate        time.Time `gorm:"not null;index" json:"match_date"`
	HomeScore        *int16    `json:"home_score"`
	AwayScore        *int16    `json:"away_score"`
	HomePenaltyScore *int16    `json:"home_penalty_score"`
	AwayPenaltyScore *int16    `json:"away_penalty_score"`
	Minute           *int16    `json:"minute"`
	Round            string    `gorm:"type:varchar(50);default:'';not null" json:"round"`
	Status           string    `gorm:"type:varchar(20);default:scheduled;not null" json:"status"`
	ExternalAPIID    *string   `gorm:"uniqueIndex" json:"-"` // id สำหรับซิงก์ external API (ไม่แสดงใน JSON)
}
