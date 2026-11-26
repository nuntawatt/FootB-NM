package models

import "github.com/google/uuid"

type Prediction struct {
	Base
	UserID             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_match" json:"user_id"`
	MatchID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_match" json:"match_id"`
	PredictedHomeScore int16     `gorm:"not null" json:"predicted_home_score"`
	PredictedAwayScore int16     `gorm:"not null" json:"predicted_away_score"`
	PredictedResult    string    `gorm:"type:varchar(10);not null" json:"predicted_result"` // home_win, draw, away_win
	IsCorrect          *bool     `json:"is_correct"`
	PointsEarned       int       `gorm:"default:0" json:"points_earned"`
}
