package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Match struct {
	Base
	LeagueID         uuid.UUID `gorm:"type:uuid;not null;index" json:"league_id"`
	League           *League   `gorm:"foreignKey:LeagueID" json:"league,omitempty"`
	HomeTeamID       uuid.UUID `gorm:"type:uuid;not null;index" json:"home_team_id"`
	HomeTeam         *Team     `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	AwayTeamID       uuid.UUID `gorm:"type:uuid;not null;index" json:"away_team_id"`
	AwayTeam         *Team     `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	MatchDate        time.Time `gorm:"not null;index" json:"match_date"`
	Year             int       `gorm:"-" json:"year"`
	Season           string    `gorm:"-" json:"season"`
	SeasonYear       int       `gorm:"-" json:"season_year"`
	HomeScore        *int16    `json:"home_score"`
	AwayScore        *int16    `json:"away_score"`
	HomePenaltyScore *int16    `json:"home_penalty_score"`
	AwayPenaltyScore *int16    `json:"away_penalty_score"`
	Minute           *int16    `json:"minute"`
	Round            string    `gorm:"type:varchar(50);default:'';not null" json:"round"`
	Status           string    `gorm:"type:varchar(20);default:scheduled;not null" json:"status"`
	ExternalAPIID    *string   `gorm:"uniqueIndex" json:"-"` // id สำหรับซิงก์ external API (ไม่แสดงใน JSON)
}

// AfterFind ถูกเรียกโดยอัตโนมัติเมื่อ GORM โหลดข้อมูลแมตช์จากฐานข้อมูล
func (m *Match) AfterFind(tx *gorm.DB) (err error) {
	m.ComputeSeasonInfo()
	return nil
}

// ComputeSeasonInfo คำนวณปี ค.ศ. และฤดูกาลจากวันแข่งขัน
func (m *Match) ComputeSeasonInfo() {
	if m.MatchDate.IsZero() {
		return
	}
	m.Year = m.MatchDate.Year()
	// ในฟุตบอลไทย ฤดูกาลเริ่มเดือน ก.ค. (7) ของปี และสิ้นสุดเดือน มิ.ย. (6) ของปีถัดไป
	if m.MatchDate.Month() >= time.July {
		m.SeasonYear = m.MatchDate.Year()
	} else {
		m.SeasonYear = m.MatchDate.Year() - 1
	}
	m.Season = fmt.Sprintf("%d/%d", m.SeasonYear, m.SeasonYear+1)
}
