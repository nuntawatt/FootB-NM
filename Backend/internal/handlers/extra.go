package handlers

import (
	"fmt"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExtraHandler struct{ db *gorm.DB }

func NewExtraHandler(db *gorm.DB) *ExtraHandler {
	return &ExtraHandler{db: db}
}

// TeamForm คืนผลงาน 5 นัดล่าสุดของทีมจากแมตช์จริงใน DB
func (h *ExtraHandler) TeamForm(c fiber.Ctx) error {
	teamID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid team id")
	}

	var matches []models.Match
	err = h.db.Preload("HomeTeam").Preload("AwayTeam").
		Where("(home_team_id = ? OR away_team_id = ?) AND status = 'finished'", teamID, teamID).
		Order("match_date desc").
		Limit(5).
		Find(&matches).Error
	if err != nil {
		return err
	}

	type FormMatch struct {
		MatchID  uuid.UUID `json:"match_id"`
		Date     string    `json:"date"`
		Opponent string    `json:"opponent"`
		Score    string    `json:"score"`
		Result   string    `json:"result"` // "W", "D", "L"
		IsHome   bool      `json:"is_home"`
	}

	var forms []FormMatch
	for _, m := range matches {
		if m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		isHome := m.HomeTeamID == teamID
		var oppName string
		var res string
		hs, as := *m.HomeScore, *m.AwayScore

		if isHome {
			oppName = m.AwayTeam.Name
			if hs > as {
				res = "W"
			} else if hs < as {
				res = "L"
			} else {
				res = "D"
			}
		} else {
			oppName = m.HomeTeam.Name
			if as > hs {
				res = "W"
			} else if as < hs {
				res = "L"
			} else {
				res = "D"
			}
		}

		forms = append(forms, FormMatch{
			MatchID:  m.ID,
			Date:     m.MatchDate.Format("2006-01-02"),
			Opponent: oppName,
			Score:    fmt.Sprintf("%d - %d", hs, as),
			Result:   res,
			IsHome:   isHome,
		})
	}

	return c.JSON(forms)
}

// TeamCoach คืนข้อมูลโค้ชของทีม (ค้นหาจาก DB ก่อน ถ้าไม่มีให้ fallback รูปและชื่อ)
func (h *ExtraHandler) TeamCoach(c fiber.Ctx) error {
	teamID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid team id")
	}

	var team models.Team
	if err := h.db.First(&team, "id = ?", teamID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "team not found")
	}

	var coach models.Coach
	if err := h.db.Where("team_id = ?", teamID).First(&coach).Error; err == nil {
		return c.JSON(coach)
	}

	// Fallback หากยังไม่ได้รัน sync
	return c.JSON(models.Coach{
		TeamID:      teamID,
		Name:        "Head Coach of " + team.Name,
		Nationality: "Thailand",
		PhotoURL:    "https://ui-avatars.com/api/?name=Head+Coach&background=0F172A&color=38BDF8&size=256",
	})
}

// TeamSquad คืนรายชื่อผู้เล่นในทีม (ค้นหาจาก DB ก่อน)
func (h *ExtraHandler) TeamSquad(c fiber.Ctx) error {
	teamID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid team id")
	}

	var players []models.Player
	if err := h.db.Where("team_id = ?", teamID).Order("number asc").Find(&players).Error; err == nil && len(players) > 0 {
		return c.JSON(players)
	}

	// Fallback ตัวอย่างหากยังไม่ได้รัน sync
	num1, num3, num9, age25, age30 := 1, 3, 9, 25, 30
	fallback := []models.Player{
		{TeamID: teamID, Name: "Siwarak Tedsungnoen", Number: &num1, Position: "Goalkeeper", Age: &age30, PhotoURL: "https://ui-avatars.com/api/?name=Siwarak&background=1e293b&color=fff"},
		{TeamID: teamID, Name: "Pansa Hemviboon", Number: &num3, Position: "Defender", Age: &age30, PhotoURL: "https://ui-avatars.com/api/?name=Pansa&background=1e293b&color=fff"},
		{TeamID: teamID, Name: "Supachai Chaided", Number: &num9, Position: "Attacker", Age: &age25, PhotoURL: "https://ui-avatars.com/api/?name=Supachai&background=1e293b&color=fff"},
	}
	return c.JSON(fallback)
}

// MatchLineups คืนรายชื่อ 11 ตัวจริงของแมตช์
func (h *ExtraHandler) MatchLineups(c fiber.Ctx) error {
	matchID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid match id")
	}

	var match models.Match
	if err := h.db.Preload("HomeTeam").Preload("AwayTeam").First(&match, "id = ?", matchID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "match not found")
	}

	num1, num2, num4, num5, num3, num6, num8, num10, num7, num9, num11 := 1, 2, 4, 5, 3, 6, 8, 10, 7, 9, 11
	lineups := models.MatchLineup{
		MatchID: match.ID,
		Home: models.LineupTeam{
			TeamID:    match.HomeTeamID,
			TeamName:  match.HomeTeam.Name,
			Formation: "4-3-3",
			CoachName: "Home Head Coach",
			Starting: []models.LineupPlayer{
				{Number: &num1, Name: "Goalkeeper", Position: "GK", Rating: "7.0"},
				{Number: &num2, Name: "Right Back", Position: "DF", Rating: "6.9"},
				{Number: &num4, Name: "Center Back", Position: "DF", Rating: "7.1"},
				{Number: &num5, Name: "Center Back", Position: "DF", Rating: "7.0"},
				{Number: &num3, Name: "Left Back", Position: "DF", Rating: "7.3"},
				{Number: &num6, Name: "Defensive Mid", Position: "MF", Rating: "6.8"},
				{Number: &num8, Name: "Central Mid", Position: "MF", Rating: "7.4"},
				{Number: &num10, Name: "Attacking Mid", Position: "MF", Rating: "7.5"},
				{Number: &num7, Name: "Right Wing", Position: "FW", Rating: "7.2"},
				{Number: &num9, Name: "Striker", Position: "FW", Rating: "8.0"},
				{Number: &num11, Name: "Left Wing", Position: "FW", Rating: "7.1"},
			},
			Subs: []models.LineupPlayer{},
		},
		Away: models.LineupTeam{
			TeamID:    match.AwayTeamID,
			TeamName:  match.AwayTeam.Name,
			Formation: "4-2-3-1",
			CoachName: "Away Head Coach",
			Starting: []models.LineupPlayer{
				{Number: &num1, Name: "Goalkeeper", Position: "GK", Rating: "6.8"},
				{Number: &num2, Name: "Right Back", Position: "DF", Rating: "6.9"},
				{Number: &num4, Name: "Center Back", Position: "DF", Rating: "7.0"},
				{Number: &num5, Name: "Center Back", Position: "DF", Rating: "7.1"},
				{Number: &num3, Name: "Left Back", Position: "DF", Rating: "6.7"},
				{Number: &num6, Name: "Central Mid", Position: "MF", Rating: "7.0"},
				{Number: &num8, Name: "Central Mid", Position: "MF", Rating: "7.2"},
				{Number: &num7, Name: "Right Wing", Position: "FW", Rating: "6.9"},
				{Number: &num10, Name: "Playmaker", Position: "MF", Rating: "7.6"},
				{Number: &num11, Name: "Left Wing", Position: "FW", Rating: "7.0"},
				{Number: &num9, Name: "Striker", Position: "FW", Rating: "7.4"},
			},
			Subs: []models.LineupPlayer{},
		},
	}
	return c.JSON(lineups)
}

// MatchVenue คืนข้อมูลสนามและช่องถ่ายทอดสด
func (h *ExtraHandler) MatchVenue(c fiber.Ctx) error {
	matchID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid match id")
	}

	var match models.Match
	if err := h.db.Preload("HomeTeam").First(&match, "id = ?", matchID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "match not found")
	}

	venue := models.VenueInfo{
		MatchID:   match.ID,
		Stadium:   match.HomeTeam.Name + " Stadium",
		City:      "Bangkok, Thailand",
		Capacity:  25000,
		Broadcast: []string{"TrueID Sports 1", "AIS PLAY Football 2", "PPTV HD 36"},
	}
	return c.JSON(venue)
}

// Heroes คืนแบนเนอร์ข่าวเด่นหน้าแรก
func (h *ExtraHandler) Heroes(c fiber.Ctx) error {
	var nextMatch models.Match
	h.db.Where("status = 'scheduled'").Order("match_date asc").First(&nextMatch)

	var matchID *uuid.UUID
	if nextMatch.ID != uuid.Nil {
		matchID = &nextMatch.ID
	}

	heroes := []models.HeroHighlight{
		{
			ID:          "hero-1",
			Title:       "ศึกบิ๊กแมตช์ไทยลีก นัดหยุดโลกสุดสัปดาห์นี้!",
			Subtitle:    "ติดตามการถ่ายทอดสดแบบ Real-time พร้อมผลการวิเคราะห์เจาะลึก",
			ImageURL:    "https://images.unsplash.com/photo-1518091043644-c1d4457512c6?auto=format&fit=crop&w=1200&q=80",
			Tag:         "BIG MATCH",
			MatchID:     matchID,
			RedirectURL: "/matches",
		},
		{
			ID:          "hero-2",
			Title:       "สรุปตารางคะแนนล่าสุด ใครเต็งแชมป์ประจำฤดูกาล?",
			Subtitle:    "เจาะลึกฟอร์ม 5 นัดหลังสุดของหัวตาราง",
			ImageURL:    "https://images.unsplash.com/photo-1579952363873-27f3bade9f55?auto=format&fit=crop&w=1200&q=80",
			Tag:         "ANALYSIS",
			RedirectURL: "/standings",
		},
	}
	return c.JSON(heroes)
}
