package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	cache "football/internal/redis"
	"football/internal/models"
	"football/internal/news"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExtraHandler struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewExtraHandler(db *gorm.DB, cache *cache.Cache) *ExtraHandler {
	return &ExtraHandler{db: db, cache: cache}
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

	var forms []models.FormMatch
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

		forms = append(forms, models.FormMatch{
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

	var teams []models.Team
	if err := h.db.Limit(1).Find(&teams, "id = ?", teamID).Error; err != nil {
		return err
	}
	if len(teams) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "team not found")
	}
	team := teams[0]

	var coaches []models.Coach
	if err := h.db.Where("team_id = ?", teamID).Limit(1).Find(&coaches).Error; err == nil && len(coaches) > 0 {
		return c.JSON(coaches[0])
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

// buildTeamLineup ดึงรายชื่อ 11 ตัวจริง โค้ช และตัวสำรองจากฐานข้อมูลจริง
func (h *ExtraHandler) buildTeamLineup(teamID uuid.UUID, teamName string, defaultFormation string) models.LineupTeam {
	// 1. ดึงข้อมูลโค้ชจริง
	coachName := "Head Coach"
	var coaches []models.Coach
	if err := h.db.Where("team_id = ?", teamID).Limit(1).Find(&coaches).Error; err == nil && len(coaches) > 0 {
		coachName = coaches[0].Name
	}

	// 2. ดึงนักเตะจริงของทีมเรียงตามเบอร์เสื้อ
	var players []models.Player
	_ = h.db.Where("team_id = ?", teamID).Order("number asc NULLS LAST").Find(&players).Error

	if len(players) == 0 {
		// Fallback ตัวอย่างหากทีมยังไม่มีข้อมูลนักเตะ
		num1, num2, num4, num5, num3, num6, num8, num10, num7, num9, num11 := 1, 2, 4, 5, 3, 6, 8, 10, 7, 9, 11
		return models.LineupTeam{
			TeamID:    teamID,
			TeamName:  teamName,
			Formation: defaultFormation,
			CoachName: coachName,
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
		}
	}

	// แยกนักเตะตามตำแหน่ง
	var gks, dfs, mfs, fws, others []models.Player
	for _, p := range players {
		pos := strings.ToLower(p.Position)
		switch {
		case strings.Contains(pos, "goal") || strings.Contains(pos, "gk"):
			gks = append(gks, p)
		case strings.Contains(pos, "def") || strings.Contains(pos, "back") || strings.Contains(pos, "df"):
			dfs = append(dfs, p)
		case strings.Contains(pos, "mid") || strings.Contains(pos, "mf"):
			mfs = append(mfs, p)
		case strings.Contains(pos, "att") || strings.Contains(pos, "forw") || strings.Contains(pos, "wing") || strings.Contains(pos, "striker") || strings.Contains(pos, "fw"):
			fws = append(fws, p)
		default:
			others = append(others, p)
		}
	}

	var starting []models.LineupPlayer
	usedIDs := make(map[uuid.UUID]bool)
	ratings := []string{"7.0", "7.2", "6.9", "7.1", "7.3", "6.8", "7.4", "7.5", "7.2", "7.8", "7.1"}

	addPlayer := func(p models.Player, posLabel string) {
		rating := "7.0"
		if len(starting) < len(ratings) {
			rating = ratings[len(starting)]
		}
		starting = append(starting, models.LineupPlayer{
			Number:   p.Number,
			Name:     p.Name,
			Position: posLabel,
			Rating:   rating,
		})
		usedIDs[p.ID] = true
	}

	// 1. ผู้รักษาประตู 1 คน
	if len(gks) > 0 {
		addPlayer(gks[0], "GK")
	}

	// 2. กองหลัง 4 คน
	for i := 0; i < len(dfs) && len(starting) < 5; i++ {
		addPlayer(dfs[i], "DF")
	}

	// 3. กองกลาง 3 คน
	for i := 0; i < len(mfs) && len(starting) < 8; i++ {
		addPlayer(mfs[i], "MF")
	}

	// 4. กองหน้า 3 คน
	for i := 0; i < len(fws) && len(starting) < 11; i++ {
		addPlayer(fws[i], "FW")
	}

	// หากตำแหน่งใดไม่พอ ให้เติมจากนักเตะที่เหลือจนครบ 11 คน (ถ้ามีพอ)
	allRemaining := append(append(append(append(gks, dfs...), mfs...), fws...), others...)
	for _, p := range allRemaining {
		if len(starting) >= 11 {
			break
		}
		if !usedIDs[p.ID] {
			posLabel := "MF"
			pos := strings.ToLower(p.Position)
			if strings.Contains(pos, "goal") || strings.Contains(pos, "gk") {
				posLabel = "GK"
			} else if strings.Contains(pos, "def") || strings.Contains(pos, "back") {
				posLabel = "DF"
			} else if strings.Contains(pos, "att") || strings.Contains(pos, "forw") {
				posLabel = "FW"
			}
			addPlayer(p, posLabel)
		}
	}

	// ผู้เล่นที่เหลือจัดเป็นตัวสำรอง (Substitutes)
	var subs []models.LineupPlayer
	for _, p := range allRemaining {
		if !usedIDs[p.ID] {
			posLabel := "MF"
			pos := strings.ToLower(p.Position)
			if strings.Contains(pos, "goal") || strings.Contains(pos, "gk") {
				posLabel = "GK"
			} else if strings.Contains(pos, "def") || strings.Contains(pos, "back") {
				posLabel = "DF"
			} else if strings.Contains(pos, "att") || strings.Contains(pos, "forw") {
				posLabel = "FW"
			}
			subs = append(subs, models.LineupPlayer{
				Number:   p.Number,
				Name:     p.Name,
				Position: posLabel,
			})
			usedIDs[p.ID] = true
		}
	}

	return models.LineupTeam{
		TeamID:    teamID,
		TeamName:  teamName,
		Formation: defaultFormation,
		CoachName: coachName,
		Starting:  starting,
		Subs:      subs,
	}
}

// MatchLineups คืนรายชื่อ 11 ตัวจริงของแมตช์
func (h *ExtraHandler) MatchLineups(c fiber.Ctx) error {
	matchID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid match id")
	}

	var matches []models.Match
	if err := h.db.Preload("HomeTeam").Preload("AwayTeam").Limit(1).Find(&matches, "id = ?", matchID).Error; err != nil {
		return err
	}
	if len(matches) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "match not found")
	}
	match := matches[0]

	lineups := models.MatchLineup{
		MatchID: match.ID,
		Home:    h.buildTeamLineup(match.HomeTeamID, match.HomeTeam.Name, "4-3-3"),
		Away:    h.buildTeamLineup(match.AwayTeamID, match.AwayTeam.Name, "4-2-3-1"),
	}
	return c.JSON(lineups)
}

// MatchVenue คืนข้อมูลสนามและช่องถ่ายทอดสด
func (h *ExtraHandler) MatchVenue(c fiber.Ctx) error {
	matchID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid match id")
	}

	var matches []models.Match
	if err := h.db.Preload("HomeTeam").Limit(1).Find(&matches, "id = ?", matchID).Error; err != nil {
		return err
	}
	if len(matches) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "match not found")
	}
	match := matches[0]

	venue := models.VenueInfo{
		MatchID:   match.ID,
		Stadium:   match.HomeTeam.Name + " Stadium",
		City:      "Bangkok, Thailand",
		Capacity:  25000,
		Broadcast: []string{"TrueID Sports 1", "AIS PLAY Football 2", "PPTV HD 36"},
	}
	return c.JSON(venue)
}

// News คืนข่าวสารฟุตบอลจริงสดๆ จาก Football News API พร้อมแปลไทยอัตโนมัติ
func (h *ExtraHandler) News(c fiber.Ctx) error {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > 20 {
		limit = 20
	}

	cacheKey := fmt.Sprintf("football:news:%d", limit)
	var cachedNews []models.NewsItem

	// 1. อ่านจาก Redis ก่อนถ้ามี
	if h.cache.Get(c.Context(), cacheKey, &cachedNews) && len(cachedNews) > 0 {
		return c.JSON(cachedNews)
	}

	// 2. ดึงข่าวจริงจากแหล่งข่าวฟุตบอลชั้นนำ (ESPN Soccer API) และแปลภาษาไทย
	items := news.FetchFootballHeroes(c.Context(), limit)
	if len(items) > 0 {
		// เซฟลง Redis แคชไว้ 30 นาที
		h.cache.Set(c.Context(), cacheKey, items, 30*time.Minute)
		return c.JSON(items)
	}

	return c.JSON(news.DefaultHeroes())
}

// Heroes รักษา backward compatibility สำหรับชื่อเดิม
func (h *ExtraHandler) Heroes(c fiber.Ctx) error {
	return h.News(c)
}
