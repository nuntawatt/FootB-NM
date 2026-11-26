package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string

	// ส่วนของ API-Football (ใช้กับคำสั่ง sync)
	APIFootballKey     string
	APIFootballBaseURL string
	Season             int // ปีของฤดูกาล เช่น 2026 สำหรับฤดูกาล 2026/27 (0 = ไม่ได้ตั้ง)
}

func Load() (*Config, error) {
	// โหลด .env ถ้ามี (บน server จริงจะใช้ env ของระบบแทน ไม่มีไฟล์ก็ไม่เป็นไร)
	_ = godotenv.Load()

	cfg := &Config{
		Port:               getenv("PORT", "3000"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		APIFootballKey:     os.Getenv("FOOTBALL_API_KEY"),
		APIFootballBaseURL: getenv("FOOTBALL_API_BASE_URL", "https://v3.football.api-sports.io"),
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	if s := os.Getenv("FOOTBALL_SEASON"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("FOOTBALL_SEASON must be a year like 2026, got %q", s)
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