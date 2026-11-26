package apifootball

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"football/internal/provider"
)

const (
	DefaultBaseURL  = "https://v3.football.api-sports.io"
	DefaultTimezone = "Asia/Bangkok"
)

var ErrRateLimited = errors.New("apifootball: rate limited (HTTP 429)")

type Client struct {
	baseURL     string
	apiKey      string
	timezone    string
	http        *http.Client
	minInterval time.Duration

	mu        sync.Mutex
	last      time.Time
	remainDay int // โควตารายวัน (-1 คือยังไม่ทราบ)
	remainMin int // โควตารายนาที (-1 คือยังไม่ทราบ)
}

type Option func(*Client)

func WithBaseURL(u string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(u, "/")
	}
}

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		c.http = h
	}
}

func WithMinInterval(d time.Duration) Option {
	return func(c *Client) {
		c.minInterval = d
	}
}

func WithTimezone(tz string) Option {
	return func(c *Client) {
		c.timezone = tz
	}
}

func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL:     DefaultBaseURL,
		apiKey:      apiKey,
		timezone:    DefaultTimezone,
		http:        &http.Client{Timeout: 20 * time.Second},
		minInterval: 6500 * time.Millisecond, // หน่วง 6.5s เพื่อไม่ให้เกิน 10 requests/นาที
		remainDay:   -1,
		remainMin:   -1,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

var _ provider.Provider = (*Client)(nil)

// Quota คืนค่าโควตาคงเหลือจาก response header
func (c *Client) Quota() (day, minute int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.remainDay, c.remainMin
}

// wait หน่วงเวลาตาม minInterval เพื่อกันยิงเกิน rate limit
func (c *Client) wait(ctx context.Context) error {
	c.mu.Lock()

	now := time.Now()
	next := now
	if !c.last.IsZero() {
		if candidate := c.last.Add(c.minInterval); candidate.After(now) {
			next = candidate
		}
	}
	c.last = next

	c.mu.Unlock()

	delay := time.Until(next)
	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *Client) recordQuota(h http.Header) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if n, err := strconv.Atoi(h.Get("x-ratelimit-requests-remaining")); err == nil {
		c.remainDay = n
	}
	if n, err := strconv.Atoi(h.Get("x-ratelimit-remaining")); err == nil {
		c.remainMin = n
	}
}

// hasErrors ตรวจสอบว่า API envelope มี error หรือไม่
func hasErrors(raw json.RawMessage) bool {
	switch strings.TrimSpace(string(raw)) {
	case "", "null", "[]", "{}":
		return false
	}
	return true
}

func (c *Client) get(ctx context.Context, path string, params url.Values, out any) error {
	for attempt := 0; attempt < 3; attempt++ {
		if err := c.wait(ctx); err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
		if err != nil {
			return err
		}
		req.URL.RawQuery = params.Encode()
		req.Header.Set("x-apisports-key", c.apiKey)

		resp, err := c.http.Do(req)
		if err != nil {
			return fmt.Errorf("apifootball %s: %w", path, err)
		}

		c.recordQuota(resp.Header)

		body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("apifootball %s: read body: %w", path, err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			// ติด rate limit รายนาที รอ 25 วินาทีแล้วลองใหม่
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(25 * time.Second):
				continue
			}
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("apifootball %s: HTTP %d: %s", path, resp.StatusCode, truncate(body, 300))
		}

		var env struct {
			Errors   json.RawMessage `json:"errors"`
			Response json.RawMessage `json:"response"`
		}
		if err := json.Unmarshal(body, &env); err != nil {
			return fmt.Errorf("apifootball %s: decode envelope: %w", path, err)
		}
		if hasErrors(env.Errors) {
			return fmt.Errorf("apifootball %s: api error: %s", path, strings.TrimSpace(string(env.Errors)))
		}

		if err := json.Unmarshal(env.Response, out); err != nil {
			return fmt.Errorf("apifootball %s: decode response: %w", path, err)
		}
		return nil
	}
	return ErrRateLimited
}

func truncate(b []byte, n int) string {
	s := strings.TrimSpace(string(b))
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

type leagueItem struct {
	League struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Logo string `json:"logo"`
	} `json:"league"`
	Country struct {
		Name string `json:"name"`
	} `json:"country"`
	Seasons []struct {
		Year    int  `json:"year"`
		Current bool `json:"current"`
	} `json:"seasons"`
}

type teamDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type fixtureItem struct {
	Fixture struct {
		ID     int    `json:"id"`
		Date   string `json:"date"`
		Status struct {
			Short   string `json:"short"`
			Elapsed *int   `json:"elapsed"`
		} `json:"status"`
	} `json:"fixture"`
	League struct {
		ID    int    `json:"id"`
		Round string `json:"round"`
	} `json:"league"`
	Teams struct {
		Home teamDTO `json:"home"`
		Away teamDTO `json:"away"`
	} `json:"teams"`
	Goals struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	} `json:"goals"`
	Score struct {
		Penalty struct {
			Home *int `json:"home"`
			Away *int `json:"away"`
		} `json:"penalty"`
	} `json:"score"`
}

// MapStatus แปลงสถานะจาก API-Football เป็นสถานะของระบบ
func MapStatus(short string) (status string, ok bool) {
	switch short {
	case "TBD", "NS":
		return provider.StatusScheduled, true
	case "1H", "HT", "2H", "ET", "BT", "P", "INT", "LIVE":
		return provider.StatusLive, true
	case "FT", "AET", "PEN", "AWD", "WO":
		return provider.StatusFinished, true
	case "PST", "SUSP", "CANC", "ABD":
		return provider.StatusPostponed, true
	}
	return "", false
}

func (it fixtureItem) toDomain() (provider.Fixture, error) {
	kickOff, err := time.Parse(time.RFC3339, it.Fixture.Date)
	if err != nil {
		return provider.Fixture{}, fmt.Errorf("fixture %d: bad date %q: %w", it.Fixture.ID, it.Fixture.Date, err)
	}

	status, _ := MapStatus(it.Fixture.Status.Short)

	f := provider.Fixture{
		ExternalID:       strconv.Itoa(it.Fixture.ID),
		LeagueExternalID: strconv.Itoa(it.League.ID),
		Round:            it.League.Round,
		KickOff:          kickOff,
		Status:           status,
		RawStatus:        it.Fixture.Status.Short,
		Home: provider.Team{
			ExternalID: strconv.Itoa(it.Teams.Home.ID),
			Name:       it.Teams.Home.Name,
			LogoURL:    it.Teams.Home.Logo,
		},
		Away: provider.Team{
			ExternalID: strconv.Itoa(it.Teams.Away.ID),
			Name:       it.Teams.Away.Name,
			LogoURL:    it.Teams.Away.Logo,
		},
		HomeScore:   it.Goals.Home,
		AwayScore:   it.Goals.Away,
		HomePenalty: it.Score.Penalty.Home,
		AwayPenalty: it.Score.Penalty.Away,
	}

	if status == provider.StatusLive {
		f.Minute = it.Fixture.Status.Elapsed
	}

	return f, nil
}

func toDomainFixtures(items []fixtureItem) ([]provider.Fixture, error) {
	out := make([]provider.Fixture, 0, len(items))
	for _, it := range items {
		f, err := it.toDomain()
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

func (c *Client) Leagues(ctx context.Context, country string) ([]provider.League, error) {
	q := url.Values{}
	if country != "" {
		q.Set("country", country)
	}

	var items []leagueItem
	if err := c.get(ctx, "/leagues", q, &items); err != nil {
		return nil, err
	}

	out := make([]provider.League, 0, len(items))
	for _, it := range items {
		l := provider.League{
			ExternalID: strconv.Itoa(it.League.ID),
			Name:       it.League.Name,
			Country:    it.Country.Name,
			LogoURL:    it.League.Logo,
			Type:       "league",
		}
		if strings.EqualFold(it.League.Type, "cup") {
			l.Type = "cup"
		}
		for _, s := range it.Seasons {
			if s.Current {
				l.CurrentSeason = s.Year
			}
		}
		out = append(out, l)
	}
	return out, nil
}

func (c *Client) Fixtures(ctx context.Context, fq provider.FixtureQuery) ([]provider.Fixture, error) {
	q := url.Values{}
	q.Set("league", fq.LeagueExternalID)
	q.Set("season", strconv.Itoa(fq.Season))
	q.Set("timezone", c.timezone)

	if !fq.From.IsZero() {
		q.Set("from", fq.From.Format("2006-01-02"))
	}
	if !fq.To.IsZero() {
		q.Set("to", fq.To.Format("2006-01-02"))
	}

	var items []fixtureItem
	if err := c.get(ctx, "/fixtures", q, &items); err != nil {
		return nil, err
	}
	return toDomainFixtures(items)
}

func (c *Client) LiveFixtures(ctx context.Context, leagueExternalIDs []string) ([]provider.Fixture, error) {
	if len(leagueExternalIDs) == 0 {
		return nil, nil
	}

	q := url.Values{}
	q.Set("live", strings.Join(leagueExternalIDs, "-"))
	q.Set("timezone", c.timezone)

	var items []fixtureItem
	if err := c.get(ctx, "/fixtures", q, &items); err != nil {
		return nil, err
	}
	return toDomainFixtures(items)
}

type coachItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
	Photo       string `json:"photo"`
}

type squadItem struct {
	Players []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Age      *int   `json:"age"`
		Number   *int   `json:"number"`
		Position string `json:"position"`
		Photo    string `json:"photo"`
	} `json:"players"`
}

func (c *Client) Coach(ctx context.Context, teamExtID string) (*provider.Coach, error) {
	q := url.Values{}
	q.Set("team", teamExtID)

	var items []coachItem
	if err := c.get(ctx, "/coachs", q, &items); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	first := items[0]
	return &provider.Coach{
		ExternalID:  strconv.Itoa(first.ID),
		Name:        first.Name,
		Nationality: first.Nationality,
		PhotoURL:    first.Photo,
	}, nil
}

func (c *Client) Squad(ctx context.Context, teamExtID string) ([]provider.Player, error) {
	q := url.Values{}
	q.Set("team", teamExtID)

	var items []squadItem
	if err := c.get(ctx, "/players/squads", q, &items); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	out := make([]provider.Player, 0, len(items[0].Players))
	for _, p := range items[0].Players {
		out = append(out, provider.Player{
			ExternalID: strconv.Itoa(p.ID),
			Name:       p.Name,
			Number:     p.Number,
			Position:   p.Position,
			Age:        p.Age,
			PhotoURL:   p.Photo,
		})
	}
	return out, nil
}

