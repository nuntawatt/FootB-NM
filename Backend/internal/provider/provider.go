// Package provider defines the domain types and interfaces used to fetch
// football data from any external source. Concrete implementations live in
// sub-packages (e.g. apifootball).
package provider

import (
	"context"
	"time"
)

// Match status constants used across the entire application.
// These are stored in the database and returned by the API.
const (
	StatusScheduled = "scheduled" // match has not started yet
	StatusLive      = "live"      // match is in progress
	StatusFinished  = "finished"  // match has ended
	StatusPostponed = "postponed" // match was delayed, suspended, or cancelled
)

// League represents a competition fetched from an external provider.
type League struct {
	ExternalID    string
	Name          string
	Country       string
	LogoURL       string
	Type          string // "league" or "cup"
	CurrentSeason int    // 0 means unknown
}

// Team represents a club or national team fetched from an external provider.
type Team struct {
	ExternalID string
	Name       string
	LogoURL    string
}

// Fixture represents a single scheduled or completed match.
type Fixture struct {
	ExternalID       string
	LeagueExternalID string
	Round            string
	KickOff          time.Time
	Status           string // one of the Status* constants, or "" if unrecognised
	RawStatus        string // the raw short code from the provider (e.g. "FT", "1H")
	Minute           *int   // elapsed minutes; non-nil only when Status == StatusLive
	Home, Away       Team
	HomeScore        *int // nil before the match starts
	AwayScore        *int
	HomePenalty      *int // nil when there was no penalty shoot-out
	AwayPenalty      *int
}

// FixtureQuery holds the filter parameters for the Fixtures method.
type FixtureQuery struct {
	LeagueExternalID string
	Season           int
	From, To         time.Time // zero value means no date filter
}

// Provider is the interface that every external data source must implement.
// All methods accept a context so callers can cancel or time out requests.
type Provider interface {
	// Leagues returns competitions for the given country (empty string = all).
	Leagues(ctx context.Context, country string) ([]League, error)

	// Fixtures returns matches that match the given query.
	Fixtures(ctx context.Context, q FixtureQuery) ([]Fixture, error)

	// LiveFixtures returns matches that are currently in progress
	// for the specified league external IDs.
	LiveFixtures(ctx context.Context, leagueExternalIDs []string) ([]Fixture, error)
}
