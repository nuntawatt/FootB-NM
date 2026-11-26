package models

import "github.com/google/uuid"

type StandingsRow struct {
	Rank         int       `json:"rank"`
	TeamID       uuid.UUID `json:"team_id"`
	TeamName     string    `json:"team_name"`
	TeamLogo     string    `json:"team_logo"`
	Played       int       `json:"played"`
	Won          int       `json:"won"`
	Drawn        int       `json:"drawn"`
	Lost         int       `json:"lost"`
	GoalsFor     int       `json:"goals_for"`
	GoalsAgainst int       `json:"goals_against"`
	GoalDiff     int       `json:"goal_diff"`
	Points       int       `json:"points"`
	Form         []string  `json:"form"`
}
