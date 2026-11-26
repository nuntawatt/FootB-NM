// Package config loads application configuration from environment variables.
// A .env file is loaded automatically when present (development convenience);
// in production, the variables should be set in the host environment directly.
package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration for the application.
type Config struct {
	// Server
	Port string // HTTP listen port; default "3000"

	// Database
	DatabaseURL string // PostgreSQL DSN; required

	// API-Football (used by the sync command and sync endpoint)
	APIFootballKey     string // FOOTBALL_API_KEY; warn on empty
	APIFootballBaseURL string // FOOTBALL_API_BASE_URL; default is the v3 endpoint
	Season             int    // FOOTBALL_SEASON; 0 means not set
}

// Load reads configuration from the environment (and .env if present).
// Returns an error only for truly unrecoverable problems (missing DATABASE_URL
// or malformed FOOTBALL_SEASON). A missing API key is logged as a warning
// because the HTTP server can still serve cached database data without it.
func Load() (*Config, error) {
	// Load .env silently; on a real server the file won't exist and that is fine.
	_ = godotenv.Load()

	cfg := &Config{
		Port:               getenv("PORT", "3000"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		APIFootballKey:     os.Getenv("FOOTBALL_API_KEY"),
		APIFootballBaseURL: getenv("FOOTBALL_API_BASE_URL", "https://v3.football.api-sports.io"),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required but not set")
	}

	if cfg.APIFootballKey == "" {
		log.Println("warning: FOOTBALL_API_KEY is not set; sync endpoints will fail")
	}

	if s := os.Getenv("FOOTBALL_SEASON"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("FOOTBALL_SEASON must be a year (e.g. 2026), got %q", s)
		}
		cfg.Season = n
	}

	return cfg, nil
}

// getenv returns the value of key, or fallback when the variable is unset or empty.
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
