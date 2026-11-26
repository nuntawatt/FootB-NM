package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // soft delete
}

type User struct {
	Base
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"` // ห้ามส่งออก JSON
	AvatarURL    string `json:"avatar_url"`
	Role         string `gorm:"type:varchar(10);default:user;not null" json:"role"`
}

type League struct {
	Base
	Name          string  `gorm:"not null" json:"name"`
	ShortName     string  `json:"short_name"`
	Country       string  `json:"country"`
	LogoURL       string  `json:"logo_url"`
	ExternalAPIID *string `gorm:"uniqueIndex" json:"-"` // ปิดบัง external_api_id ไม่ให้ส่งออกทาง API
	Type          string  `gorm:"type:varchar(10);default:league;not null" json:"type"`
	SortOrder     int     `gorm:"default:0;not null" json:"sort_order"`
}

type Team struct {
	Base
	LeagueID      uuid.UUID `gorm:"type:uuid;not null;index" json:"league_id"`
	Name          string    `gorm:"not null" json:"name"`
	ShortName     string    `json:"short_name"`
	LogoURL       string    `json:"logo_url"`
	ExternalAPIID *string   `gorm:"uniqueIndex" json:"-"` // ปิดบัง external_api_id ไม่ให้ส่งออกทาง API
}

type Match struct {
	Base
	LeagueID         uuid.UUID `gorm:"type:uuid;not null;index" json:"league_id"`
	HomeTeamID       uuid.UUID `gorm:"type:uuid;not null;index" json:"home_team_id"`
	AwayTeamID       uuid.UUID `gorm:"type:uuid;not null;index" json:"away_team_id"`
	MatchDate        time.Time `gorm:"not null;index" json:"match_date"`
	HomeScore        *int16    `json:"home_score"`
	AwayScore        *int16    `json:"away_score"`
	HomePenaltyScore *int16    `json:"home_penalty_score"`
	AwayPenaltyScore *int16    `json:"away_penalty_score"`
	Minute           *int16    `json:"minute"`
	Round            string    `gorm:"type:varchar(50);default:'';not null" json:"round"`
	Status           string    `gorm:"type:varchar(20);default:scheduled;not null" json:"status"`
	ExternalAPIID    *string   `gorm:"uniqueIndex" json:"-"` // ปิดบัง external_api_id ไม่ให้ส่งออกทาง API
}

type Prediction struct {
	Base
	UserID             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_match" json:"user_id"`
	MatchID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_match" json:"match_id"`
	PredictedHomeScore int16     `gorm:"not null" json:"predicted_home_score"`
	PredictedAwayScore int16     `gorm:"not null" json:"predicted_away_score"`
	PredictedResult    string    `gorm:"type:varchar(10);not null" json:"predicted_result"` // home_win | draw | away_win
	IsCorrect          *bool     `json:"is_correct"`
	PointsEarned       int       `gorm:"default:0" json:"points_earned"`
}

type UserFavoriteTeam struct {
	UserID    uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"user_id"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"team_id"`
	CreatedAt time.Time `json:"created_at"`
}