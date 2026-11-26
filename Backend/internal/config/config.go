package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	APIFootballKey     string
	APIFootballBaseURL string
	Season             int
}

func Load() (*Config, error) {
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
		log.Println("warning: FOOTBALL_API_KEY is not set")
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

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
