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
	Type          string // "league" | "cup"
	CurrentSeason int    // 0 = ไม่ทราบ
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
	Status           string // Status หรือ "" ถ้าไม่รู้จักรหัสสถานะ
	RawStatus        string
	Minute           *int // มีค่าเฉพาะตอน Match
	Home, Away       Team
	HomeScore        *int // nil = ยังไม่มีสกอร์
	AwayScore        *int
	HomePenalty      *int // คะแนนดวลจุดโทษ nil ไม่มีจุดโทษ
	AwayPenalty      *int
}

type FixtureQuery struct {
	LeagueExternalID string
	Season           int
	From, To         time.Time // ตามโซนเวลาค่าที่ส่งมา
}

type Provider interface {
	Leagues(ctx context.Context, country string) ([]League, error)
	Fixtures(ctx context.Context, q FixtureQuery) ([]Fixture, error)
	
	// LiveFixtures คืนแมตซ์ที่กำลังแข่งของลีกที่ระบุ
	LiveFixtures(ctx context.Context, leagueExternalIDs []string) ([]Fixture, error)
}
