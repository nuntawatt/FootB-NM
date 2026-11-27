package news

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestCleanHTML(t *testing.T) {
	input := "<p>&#8220;สวายเรียง&#8221; <strong>สร้างประวัติศาสตร์</strong></p>"
	expected := "“สวายเรียง” สร้างประวัติศาสตร์"
	got := cleanHTML(input)
	if got != expected {
		t.Errorf("cleanHTML(%q) = %q, expected %q", input, got, expected)
	}
}

func TestFetchFootballHeroesThai(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	heroes := FetchFootballHeroes(ctx, 10)
	if len(heroes) == 0 {
		t.Fatal("expected at least 1 hero highlight")
	}

	if len(heroes) != 10 {
		t.Errorf("expected 10 heroes, got %d", len(heroes))
	}

	t.Logf("Successfully fetched %d Thai football news items:", len(heroes))
	for i, h := range heroes {
		if h.ID == "" || h.Title == "" || h.Subtitle == "" || h.ImageURL == "" || h.Tag == "" {
			t.Errorf("missing required field in hero[%d]: %+v", i, h)
		}
		if strings.Contains(h.Title, "<") || strings.Contains(h.Title, ">") || strings.Contains(h.Title, "&#") {
			t.Errorf("title contains uncleaned html entities: %s", h.Title)
		}
		if strings.Contains(h.Subtitle, "<") || strings.Contains(h.Subtitle, ">") || strings.Contains(h.Subtitle, "&#") {
			t.Errorf("subtitle contains uncleaned html entities: %s", h.Subtitle)
		}
		t.Logf("[%d] ID=%s | Tag=%s\nTitle: %s\nSubtitle: %s\nImage: %s\nLink: %s\n",
			i+1, h.ID, h.Tag, h.Title, h.Subtitle, h.ImageURL, h.RedirectURL)
	}
}
