package provider

import (
	"context"
	"time"
)

const (
	StatusScheduled = "scheduled"
	StatusLive      = "live"
	StatusFinished  = "finished"
	StatusPostponed = "postponed"
)

type League struct {
	ExternalID    string
	Name          string
	Country       string
	LogoURL       string
	Type          string // "league" หรือ "cup"
	CurrentSeason int
}

type Team struct {
	ExternalID string
	Name       string
	LogoURL    string
}

type Fixture struct {
	ExternalID       string
	LeagueExternalID string
	Round            string
	KickOff          time.Time
	Status           string
	RawStatus        string
	Minute           *int
	Home, Away       Team
	HomeScore        *int
	AwayScore        *int
	HomePenalty      *int
	AwayPenalty      *int
}

type FixtureQuery struct {
	LeagueExternalID string
	Season           int
	From, To         time.Time
}

type Provider interface {
	Leagues(ctx context.Context, country string) ([]League, error)
	Fixtures(ctx context.Context, q FixtureQuery) ([]Fixture, error)
	LiveFixtures(ctx context.Context, leagueExternalIDs []string) ([]Fixture, error)
}
