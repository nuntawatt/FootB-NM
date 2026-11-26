package models

import "github.com/google/uuid"

type LineupPlayer struct {
	Number   *int   `json:"number"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Rating   string `json:"rating,omitempty"`
}

type LineupTeam struct {
	TeamID    uuid.UUID      `json:"team_id"`
	TeamName  string         `json:"team_name"`
	Formation string         `json:"formation"`
	CoachName string         `json:"coach_name"`
	Starting  []LineupPlayer `json:"starting"`
	Subs      []LineupPlayer `json:"substitutes"`
}

type MatchLineup struct {
	MatchID uuid.UUID  `json:"match_id"`
	Home    LineupTeam `json:"home"`
	Away    LineupTeam `json:"away"`
}

type VenueInfo struct {
	MatchID   uuid.UUID `json:"match_id"`
	Stadium   string    `json:"stadium"`
	City      string    `json:"city"`
	Capacity  int       `json:"capacity"`
	Broadcast []string  `json:"broadcast"`
}

type HeroHighlight struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Subtitle    string     `json:"subtitle"`
	ImageURL    string     `json:"image_url"`
	Tag         string     `json:"tag"`
	MatchID     *uuid.UUID `json:"match_id,omitempty"`
	RedirectURL string     `json:"redirect_url,omitempty"`
}

type FormMatch struct {
	MatchID  uuid.UUID `json:"match_id"`
	Date     string    `json:"date"`
	Opponent string    `json:"opponent"`
	Score    string    `json:"score"`
	Result   string    `json:"result"` // "W", "D", "L"
	IsHome   bool      `json:"is_home"`
}

