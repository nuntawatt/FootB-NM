package models

import (
	"testing"
	"time"
)

func TestMatchComputeSeasonInfo(t *testing.T) {
	tests := []struct {
		name       string
		matchDate  time.Time
		wantYear   int
		wantSeason string
		wantSYear  int
	}{
		{
			name:       "Late 2024 (Season 2024/2025)",
			matchDate:  time.Date(2024, time.August, 10, 18, 0, 0, 0, time.UTC),
			wantYear:   2024,
			wantSeason: "2024/2025",
			wantSYear:  2024,
		},
		{
			name:       "Early 2025 (Season 2024/2025)",
			matchDate:  time.Date(2025, time.May, 24, 18, 0, 0, 0, time.UTC),
			wantYear:   2025,
			wantSeason: "2024/2025",
			wantSYear:  2024,
		},
		{
			name:       "Late 2022 (Season 2022/2023)",
			matchDate:  time.Date(2022, time.November, 27, 19, 0, 0, 0, time.UTC),
			wantYear:   2022,
			wantSeason: "2022/2023",
			wantSYear:  2022,
		},
		{
			name:       "Early 2023 (Season 2022/2023)",
			matchDate:  time.Date(2023, time.March, 15, 17, 30, 0, 0, time.UTC),
			wantYear:   2023,
			wantSeason: "2022/2023",
			wantSYear:  2022,
		},
		{
			name:       "Late 2025 (Season 2025/2026)",
			matchDate:  time.Date(2025, time.August, 1, 19, 0, 0, 0, time.UTC),
			wantYear:   2025,
			wantSeason: "2025/2026",
			wantSYear:  2025,
		},
		{
			name:       "Early 2026 (Season 2025/2026)",
			matchDate:  time.Date(2026, time.February, 20, 18, 0, 0, 0, time.UTC),
			wantYear:   2026,
			wantSeason: "2025/2026",
			wantSYear:  2025,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := Match{MatchDate: tc.matchDate}
			m.ComputeSeasonInfo()

			if m.Year != tc.wantYear {
				t.Errorf("Year = %d, want %d", m.Year, tc.wantYear)
			}
			if m.Season != tc.wantSeason {
				t.Errorf("Season = %q, want %q", m.Season, tc.wantSeason)
			}
			if m.SeasonYear != tc.wantSYear {
				t.Errorf("SeasonYear = %d, want %d", m.SeasonYear, tc.wantSYear)
			}
		})
	}
}
