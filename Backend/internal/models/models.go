// Package models defines GORM database models for the football application.
// All models embed Base, which provides UUID primary key, timestamps, and soft delete.
//
// Design notes:
//   - ExternalAPIID is hidden from JSON output (json:"-") to avoid exposing
//     third-party identifiers in the public API.
//   - Scores and time fields use pointers (*int16, *time.Time) so that a nil
//     value means "not yet available", distinct from an explicit 0.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base is embedded by every model and provides standard audit fields.
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                           json:"-"` // soft delete; excluded from JSON
}

// User represents an application user account.
type User struct {
	Base
	Username     string `gorm:"uniqueIndex;not null"                      json:"username"`
	Email        string `gorm:"uniqueIndex;not null"                      json:"email"`
	PasswordHash string `gorm:"not null"                                  json:"-"`          // never exposed in JSON
	AvatarURL    string `                                                 json:"avatar_url"`
	Role         string `gorm:"type:varchar(10);default:user;not null"    json:"role"`
}

// League represents a football competition (e.g. Thai League 1, FA Cup).
type League struct {
	Base
	Name          string  `gorm:"not null"                               json:"name"`
	ShortName     string  `                                              json:"short_name"`
	Country       string  `                                              json:"country"`
	LogoURL       string  `                                              json:"logo_url"`
	Type          string  `gorm:"type:varchar(10);default:league;not null" json:"type"`       // "league" or "cup"
	SortOrder     int     `gorm:"default:0;not null"                     json:"sort_order"`
	ExternalAPIID *string `gorm:"uniqueIndex"                            json:"-"`            // internal idempotency key; not exposed
}

// Team represents a football club or squad.
// LeagueID records the league the team was first imported from; a team may
// appear in multiple competitions without changing this field.
type Team struct {
	Base
	LeagueID      uuid.UUID `gorm:"type:uuid;not null;index"  json:"league_id"`
	League        *League   `gorm:"foreignKey:LeagueID"       json:"league,omitempty"`
	Name          string    `gorm:"not null"                  json:"name"`
	ShortName     string    `                                 json:"short_name"`
	LogoURL       string    `                                 json:"logo_url"`
	ExternalAPIID *string   `gorm:"uniqueIndex"               json:"-"` // internal idempotency key; not exposed
}

// Match represents a single football fixture (scheduled or completed).
type Match struct {
	Base
	LeagueID         uuid.UUID `gorm:"type:uuid;not null;index"               json:"league_id"`
	League           *League   `gorm:"foreignKey:LeagueID"                    json:"league,omitempty"`
	HomeTeamID       uuid.UUID `gorm:"type:uuid;not null;index"               json:"home_team_id"`
	HomeTeam         *Team     `gorm:"foreignKey:HomeTeamID"                  json:"home_team,omitempty"`
	AwayTeamID       uuid.UUID `gorm:"type:uuid;not null;index"               json:"away_team_id"`
	AwayTeam         *Team     `gorm:"foreignKey:AwayTeamID"                  json:"away_team,omitempty"`
	MatchDate        time.Time `gorm:"not null;index"                         json:"match_date"`
	HomeScore        *int16    `                                               json:"home_score"`         // nil before match starts
	AwayScore        *int16    `                                               json:"away_score"`
	HomePenaltyScore *int16    `                                               json:"home_penalty_score"` // nil when no shoot-out
	AwayPenaltyScore *int16    `                                               json:"away_penalty_score"`
	Minute           *int16    `                                               json:"minute"`             // elapsed minutes; non-nil only when live
	Round            string    `gorm:"type:varchar(50);default:'';not null"   json:"round"`
	Status           string    `gorm:"type:varchar(20);default:scheduled;not null" json:"status"`
	ExternalAPIID    *string   `gorm:"uniqueIndex"                            json:"-"` // internal idempotency key; not exposed
}

// Prediction stores a user's pre-match score prediction.
type Prediction struct {
	Base
	UserID             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_match" json:"user_id"`
	MatchID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_match" json:"match_id"`
	PredictedHomeScore int16     `gorm:"not null"                                      json:"predicted_home_score"`
	PredictedAwayScore int16     `gorm:"not null"                                      json:"predicted_away_score"`
	PredictedResult    string    `gorm:"type:varchar(10);not null"                     json:"predicted_result"` // "home_win" | "draw" | "away_win"
	IsCorrect          *bool     `                                                      json:"is_correct"`
	PointsEarned       int       `gorm:"default:0"                                     json:"points_earned"`
}

// UserFavoriteTeam is a many-to-many join record between users and their
// favourite clubs. Routes for this model are not yet implemented.
type UserFavoriteTeam struct {
	UserID    uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"user_id"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"team_id"`
	CreatedAt time.Time `                                     json:"created_at"`
}
