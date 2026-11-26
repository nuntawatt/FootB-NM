package handlers

import (
	"sort"
	"strconv"
	"time"

	"football/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StandingsHandler struct{ db *gorm.DB }

func NewStandingsHandler(db *gorm.DB) *StandingsHandler {
	return &StandingsHandler{db: db}
}

// Get คำนวณตารางคะแนนจากผลการแข่งขันในฐานข้อมูล
func (h *StandingsHandler) Get(c fiber.Ctx) error {
	seasonYear := 2024
	if s := c.Query("season"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 2000 {
			seasonYear = n
		}
	}

	start := time.Date(seasonYear, 7, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(seasonYear+1, 7, 1, 0, 0, 0, 0, time.UTC)

	var teams []models.Team
	tQuery := h.db.Model(&models.Team{})
	if lid := c.Query("league_id"); lid != "" {
		if parsedLid, err := uuid.Parse(lid); err == nil {
			tQuery = tQuery.Where("league_id = ?", parsedLid)
		}
	}
	if err := tQuery.Find(&teams).Error; err != nil {
		return err
	}

	var matches []models.Match
	mQuery := h.db.Where("match_date >= ? AND match_date < ? AND status = 'finished'", start, end)
	if lid := c.Query("league_id"); lid != "" {
		if parsedLid, err := uuid.Parse(lid); err == nil {
			mQuery = mQuery.Where("league_id = ?", parsedLid)
		}
	}
	if err := mQuery.Order("match_date asc").Find(&matches).Error; err != nil {
		return err
	}

	table := make(map[uuid.UUID]*models.StandingsRow)
	for _, t := range teams {
		table[t.ID] = &models.StandingsRow{
			TeamID:   t.ID,
			TeamName: t.Name,
			TeamLogo: t.LogoURL,
			Form:     make([]string, 0),
		}
	}

	for _, m := range matches {
		if m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		hs, as := int(*m.HomeScore), int(*m.AwayScore)

		hRow, okH := table[m.HomeTeamID]
		aRow, okA := table[m.AwayTeamID]
		if !okH || !okA {
			continue
		}

		hRow.Played++
		aRow.Played++
		hRow.GoalsFor += hs
		hRow.GoalsAgainst += as
		aRow.GoalsFor += as
		aRow.GoalsAgainst += hs

		if hs > as {
			hRow.Won++
			hRow.Points += 3
			hRow.Form = append(hRow.Form, "W")
			aRow.Lost++
			aRow.Form = append(aRow.Form, "L")
		} else if hs < as {
			aRow.Won++
			aRow.Points += 3
			aRow.Form = append(aRow.Form, "W")
			hRow.Lost++
			hRow.Form = append(hRow.Form, "L")
		} else {
			hRow.Drawn++
			hRow.Points += 1
			hRow.Form = append(hRow.Form, "D")
			aRow.Drawn++
			aRow.Points += 1
			aRow.Form = append(aRow.Form, "D")
		}
	}

	var rows []models.StandingsRow
	for _, row := range table {
		if row.Played == 0 && len(matches) > 0 {
			continue
		}
		row.GoalDiff = row.GoalsFor - row.GoalsAgainst
		if len(row.Form) > 5 {
			row.Form = row.Form[len(row.Form)-5:]
		}
		rows = append(rows, *row)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Points != rows[j].Points {
			return rows[i].Points > rows[j].Points
		}
		if rows[i].GoalDiff != rows[j].GoalDiff {
			return rows[i].GoalDiff > rows[j].GoalDiff
		}
		if rows[i].GoalsFor != rows[j].GoalsFor {
			return rows[i].GoalsFor > rows[j].GoalsFor
		}
		return rows[i].TeamName < rows[j].TeamName
	})

	for i := range rows {
		rows[i].Rank = i + 1
	}

	return c.JSON(rows)
}
