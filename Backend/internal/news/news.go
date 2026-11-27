package news

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"football/internal/models"
)

type espnResponse struct {
	Header   string        `json:"header"`
	Articles []espnArticle `json:"articles"`
}

type espnArticle struct {
	ID          int            `json:"id"`
	Headline    string         `json:"headline"`
	Description string         `json:"description"`
	Published   string         `json:"published"`
	Images      []espnImage    `json:"images"`
	Links       espnLinks      `json:"links"`
	Categories  []espnCategory `json:"categories"`
}

type espnImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type espnLinks struct {
	Web struct {
		HRef string `json:"href"`
	} `json:"web"`
}

type espnCategory struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}

var defaultBannerImages = []string{
	"https://images.unsplash.com/photo-1508098682722-e99c43a406b2?auto=format&fit=crop&w=1200&q=80",
	"https://images.unsplash.com/photo-1518091043644-c1d4457512c6?auto=format&fit=crop&w=1200&q=80",
	"https://images.unsplash.com/photo-1579952363873-27f3bade9f55?auto=format&fit=crop&w=1200&q=80",
	"https://images.unsplash.com/photo-1522778119026-d647f0596c20?auto=format&fit=crop&w=1200&q=80",
}

// FetchFootballHeroes ดึงข่าวฟุตบอลจริงสดๆ จาก Football News API (ESPN)
func FetchFootballHeroes(ctx context.Context, limit int) []models.HeroHighlight {
	if limit <= 0 || limit > 10 {
		limit = 5
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// 1. ดึงข่าวจาก ESPN Premier League News API
	apiURL := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/soccer/eng.1/news?limit=%d", limit)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return DefaultHeroes()
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return DefaultHeroes()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DefaultHeroes()
	}

	var data espnResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || len(data.Articles) == 0 {
		return DefaultHeroes()
	}

	var heroes []models.HeroHighlight
	for i, a := range data.Articles {
		title := strings.TrimSpace(a.Headline)
		if title == "" {
			continue
		}

		desc := strings.TrimSpace(a.Description)
		if len([]rune(desc)) > 140 {
			desc = string([]rune(desc)[:140]) + "..."
		}

		imgURL := ""
		if len(a.Images) > 0 && a.Images[0].URL != "" {
			imgURL = a.Images[0].URL
		} else {
			imgURL = defaultBannerImages[i%len(defaultBannerImages)]
		}

		tag := "FOOTBALL"
		for _, cat := range a.Categories {
			if cat.Type == "league" && cat.Description != "" {
				tag = cat.Description
				break
			}
		}
		if tag == "FOOTBALL" && len(a.Categories) > 0 && a.Categories[0].Description != "" {
			tag = a.Categories[0].Description
		}

		link := a.Links.Web.HRef
		if link == "" {
			link = "/matches"
		}

		heroes = append(heroes, models.HeroHighlight{
			ID:          fmt.Sprintf("hero-%d", a.ID),
			Title:       title,
			Subtitle:    desc,
			ImageURL:    imgURL,
			Tag:         strings.ToUpper(tag),
			RedirectURL: link,
		})

		if len(heroes) >= limit {
			break
		}
	}

	if len(heroes) == 0 {
		return DefaultHeroes()
	}

	return heroes
}

// DefaultHeroes รายการไฮไลท์สำรองกรณีเชื่อมต่อ API ไม่สำเร็จ
func DefaultHeroes() []models.HeroHighlight {
	return []models.HeroHighlight{
		{
			ID:          "hero-1",
			Title:       "เกาะติดศึกการแข่งขันฟุตบอลและผลบอลสดประจำสัปดาห์",
			Subtitle:    "ติดตามตารางคะแนน ผลบอลสด และสถิติการแข่งขันครบทุกสโมสร",
			ImageURL:    defaultBannerImages[0],
			Tag:         "FOOTBALL NEWS",
			RedirectURL: "/matches",
		},
		{
			ID:          "hero-2",
			Title:       "สรุปตารางคะแนนล่าสุด เช็คอันดับจ่าฝูงและโซนตกชั้น",
			Subtitle:    "คำนวณแต้มและผลต่างประตูได้เสียอย่างแม่นยำด้วยระบบวิเคราะห์ข้อมูล",
			ImageURL:    defaultBannerImages[2],
			Tag:         "STANDINGS",
			RedirectURL: "/standings",
		},
	}
}
