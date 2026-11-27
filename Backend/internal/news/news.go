package news

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"football/internal/models"
)

type ballthaiPost struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Link    string `json:"link"`
	Title   struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Excerpt struct {
		Rendered string `json:"rendered"`
	} `json:"excerpt"`
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
	Embedded struct {
		FeaturedMedia []struct {
			SourceURL string `json:"source_url"`
		} `json:"wp:featuredmedia"`
		Terms [][]struct {
			Name string `json:"name"`
		} `json:"wp:term"`
	} `json:"_embedded"`
}

var defaultBannerImages = []string{
	"https://images.unsplash.com/photo-1508098682722-e99c43a406b2?auto=format&fit=crop&w=1200&q=80",
	"https://images.unsplash.com/photo-1518091043644-c1d4457512c6?auto=format&fit=crop&w=1200&q=80",
	"https://images.unsplash.com/photo-1579952363873-27f3bade9f55?auto=format&fit=crop&w=1200&q=80",
	"https://images.unsplash.com/photo-1522778119026-d647f0596c20?auto=format&fit=crop&w=1200&q=80",
}

// FetchFootballHeroes ดึงข่าวฟุตบอลจริงสดๆ จากสำนักข่าวฟุตบอลไทย (Ballthai API)
func FetchFootballHeroes(ctx context.Context, limit int) []models.HeroHighlight {
	if limit <= 0 {
		limit = 10
	}
	if limit > 20 {
		limit = 20
	}

	// 1. ดึงข่าวจากสำนักข่าวฟุตบอลไทยชั้นนำ
	heroes, err := fetchBallthaiNews(ctx, limit)
	if err == nil && len(heroes) > 0 {
		return heroes
	}

	// 2. กรณีฉุกเฉินเชื่อมต่อไม่ได้ ให้คืนค่าข้อมูลสำรอง
	return DefaultHeroes()
}

// fetchBallthaiNews ดึงข่าวสารฟุตบอลไทยสดๆ จาก Ballthai REST API
func fetchBallthaiNews(ctx context.Context, limit int) ([]models.HeroHighlight, error) {
	client := &http.Client{Timeout: 6 * time.Second}
	apiURL := fmt.Sprintf("https://www.ballthai.com/wp-json/wp/v2/posts?_embed&per_page=%d", limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ballthai api returned status: %d", resp.StatusCode)
	}

	var posts []ballthaiPost
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("no posts found")
	}

	var heroes []models.HeroHighlight
	for i, p := range posts {
		title := cleanHTML(p.Title.Rendered)
		if title == "" {
			continue
		}

		subtitle := extractSubtitle(p.Content.Rendered, p.Excerpt.Rendered)

		imageURL := ""
		if len(p.Embedded.FeaturedMedia) > 0 && p.Embedded.FeaturedMedia[0].SourceURL != "" {
			imageURL = p.Embedded.FeaturedMedia[0].SourceURL
		} else {
			imageURL = defaultBannerImages[i%len(defaultBannerImages)]
		}

		tag := "ข่าวฟุตบอล"
		if len(p.Embedded.Terms) > 0 && len(p.Embedded.Terms[0]) > 0 && p.Embedded.Terms[0][0].Name != "" {
			tag = cleanHTML(p.Embedded.Terms[0][0].Name)
		}

		link := p.Link
		if link == "" {
			link = "/matches"
		}

		heroes = append(heroes, models.HeroHighlight{
			ID:          fmt.Sprintf("hero-%d", p.ID),
			Title:       title,
			Subtitle:    subtitle,
			ImageURL:    imageURL,
			Tag:         tag,
			RedirectURL: link,
		})

		if len(heroes) >= limit {
			break
		}
	}

	return heroes, nil
}

// cleanHTML ถอด HTML Tags และแปลง HTML Entities ให้เป็นข้อความภาษาไทยบริสุทธิ์
func cleanHTML(s string) string {
	s = html.UnescapeString(s)
	var b strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			b.WriteRune(r)
		}
	}
	fields := strings.Fields(b.String())
	return strings.TrimSpace(strings.Join(fields, " "))
}

// extractSubtitle ดึงข้อความเกริ่นนำข่าวความยาวกำลังพอดีสำหรับแสดงผลบนการ์ดข่าว
func extractSubtitle(contentHTML, excerptHTML string) string {
	cleanContent := cleanHTML(contentHTML)
	if cleanContent != "" {
		runes := []rune(cleanContent)
		if len(runes) > 160 {
			return string(runes[:160]) + "..."
		}
		return cleanContent
	}

	cleanExcerpt := cleanHTML(excerptHTML)
	if cleanExcerpt != "" {
		runes := []rune(cleanExcerpt)
		if len(runes) > 160 {
			return string(runes[:160]) + "..."
		}
		return cleanExcerpt
	}

	return "ติดตามความเคลื่อนไหวและรายละเอียดการแข่งขันฟุตบอลล่าสุดได้ที่นี่"
}

// DefaultHeroes รายการไฮไลท์สำรองกรณีฉุกเฉิน
func DefaultHeroes() []models.HeroHighlight {
	return []models.HeroHighlight{
		{
			ID:          "hero-1",
			Title:       "เกาะติดศึกการแข่งขันฟุตบอลไทยและผลบอลสดประจำสัปดาห์",
			Subtitle:    "ติดตามตารางคะแนน ผลบอลสด และสถิติการแข่งขันไทยลีกครบทุกสโมสร",
			ImageURL:    defaultBannerImages[0],
			Tag:         "ไทยลีก",
			RedirectURL: "/matches",
		},
		{
			ID:          "hero-2",
			Title:       "สรุปตารางคะแนนล่าสุด เช็คอันดับจ่าฝูงและโซนตกชั้น",
			Subtitle:    "คำนวณแต้มและผลต่างประตูได้เสียอย่างแม่นยำด้วยระบบวิเคราะห์ข้อมูล",
			ImageURL:    defaultBannerImages[2],
			Tag:         "ตารางคะแนน",
			RedirectURL: "/standings",
		},
	}
}
