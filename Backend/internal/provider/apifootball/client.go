// Package apifootball implements the provider.Provider interface using the
// API-Football v3 REST API (https://v3.football.api-sports.io).
//
// Rate limiting:
//   - The free plan allows 100 requests per day and 10 per minute.
//   - The client enforces a minimum interval between requests (default 1 s)
//     to avoid being blocked. All quota figures are read from response headers
//     and can be inspected via Client.Quota().
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

// DefaultBaseURL is the production endpoint for API-Football v3.
const DefaultBaseURL = "https://v3.football.api-sports.io"

// DefaultTimezone is used when building fixture queries.
// Thai League matches are stored in Asia/Bangkok time.
const DefaultTimezone = "Asia/Bangkok"

// ErrRateLimited is returned when the server responds with HTTP 429.
var ErrRateLimited = errors.New("apifootball: rate limited (HTTP 429)")

// ----------------------------------------------------------------------------
// Client
// ----------------------------------------------------------------------------

// Client is a thread-safe HTTP client for the API-Football v3 service.
// Use New to construct one; the zero value is not valid.
type Client struct {
	baseURL     string
	apiKey      string
	timezone    string
	http        *http.Client
	minInterval time.Duration // minimum delay between consecutive requests

	mu        sync.Mutex
	last      time.Time // timestamp of the last request slot reserved
	remainDay int       // daily quota remaining (-1 = not yet known)
	remainMin int       // per-minute quota remaining (-1 = not yet known)
}

// Option is a functional option for New.
type Option func(*Client)

// WithBaseURL overrides the API base URL (useful for testing with a stub server).
func WithBaseURL(u string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(u, "/")
	}
}

// WithHTTPClient replaces the default http.Client (useful for testing).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		c.http = h
	}
}

// WithMinInterval overrides the minimum delay between requests.
// Default is 1 second, which satisfies the 10 req/min limit with headroom.
func WithMinInterval(d time.Duration) Option {
	return func(c *Client) {
		c.minInterval = d
	}
}

// WithTimezone sets the IANA timezone sent with fixture queries.
func WithTimezone(tz string) Option {
	return func(c *Client) {
		c.timezone = tz
	}
}

// New creates a Client with the given API key and optional overrides.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL:     DefaultBaseURL,
		apiKey:      apiKey,
		timezone:    DefaultTimezone,
		http:        &http.Client{Timeout: 20 * time.Second},
		minInterval: time.Second,
		remainDay:   -1,
		remainMin:   -1,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Compile-time proof that *Client satisfies provider.Provider.
var _ provider.Provider = (*Client)(nil)

// Quota returns the remaining daily and per-minute quota from the most recent
// response headers. Both values are -1 until the first successful request.
func (c *Client) Quota() (day, minute int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.remainDay, c.remainMin
}

// ----------------------------------------------------------------------------
// Internal helpers
// ----------------------------------------------------------------------------

// wait blocks until it is safe to send the next request.
// It reserves a slot in the rate-limit window atomically and then sleeps
// outside the mutex so other goroutines can queue their own slots.
func (c *Client) wait(ctx context.Context) error {
	c.mu.Lock()

	now := time.Now()
	next := now
	if !c.last.IsZero() {
		if candidate := c.last.Add(c.minInterval); candidate.After(now) {
			next = candidate
		}
	}
	c.last = next // reserve this slot

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

// recordQuota reads quota headers from the response and stores them.
//
// Header mapping:
//   - x-ratelimit-requests-remaining -> daily quota
//   - x-ratelimit-remaining          -> per-minute quota
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

// hasErrors reports whether the API envelope contains error details.
//
// API-Football always returns HTTP 200, even for authentication failures.
// A successful response has "errors": [] while a failed one has an object
// such as {"token": "Error/Missing application key"}.
func hasErrors(raw json.RawMessage) bool {
	switch strings.TrimSpace(string(raw)) {
	case "", "null", "[]", "{}":
		return false
	}
	return true
}

// get executes a single authenticated GET request and decodes the
// "response" field of the API envelope into out.
func (c *Client) get(ctx context.Context, path string, params url.Values, out any) error {
	// 1. Respect rate limit: wait for the next available slot.
	if err := c.wait(ctx); err != nil {
		return err
	}

	// 2. Build the request with the API key header.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("x-apisports-key", c.apiKey)

	// 3. Execute the request.
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("apifootball %s: %w", path, err)
	}
	defer resp.Body.Close()

	// 4. Record remaining quota from response headers before reading the body.
	c.recordQuota(resp.Header)

	// 5. Read the body (cap at 8 MiB to protect memory).
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return fmt.Errorf("apifootball %s: read body: %w", path, err)
	}

	// 6. Handle non-200 status codes.
	if resp.StatusCode == http.StatusTooManyRequests {
		return ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("apifootball %s: HTTP %d: %s", path, resp.StatusCode, truncate(body, 300))
	}

	// 7. Decode the outer envelope { "errors": ..., "response": ... }.
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

	// 8. Decode the "response" field into the caller-supplied type.
	if err := json.Unmarshal(env.Response, out); err != nil {
		return fmt.Errorf("apifootball %s: decode response: %w", path, err)
	}
	return nil
}

// truncate shortens a byte slice to at most n characters for safe logging.
func truncate(b []byte, n int) string {
	s := strings.TrimSpace(string(b))
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

// ----------------------------------------------------------------------------
// API-Football response DTOs (only the fields we actually use)
// ----------------------------------------------------------------------------

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
		Date   string `json:"date"` // RFC 3339
		Status struct {
			Short   string `json:"short"`
			Elapsed *int   `json:"elapsed"` // minutes played; nil when not live
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
		Home *int `json:"home"` // nil before the match starts
		Away *int `json:"away"`
	} `json:"goals"`
	Score struct {
		Penalty struct {
			Home *int `json:"home"` // nil when there was no shoot-out
			Away *int `json:"away"`
		} `json:"penalty"`
	} `json:"score"`
}

// ----------------------------------------------------------------------------
// Status mapping
// ----------------------------------------------------------------------------

// MapStatus converts an API-Football status short code to one of the
// application's Status* constants.
//
// Mapping table:
//
//	TBD, NS                        -> scheduled
//	1H, HT, 2H, ET, BT, P, INT, LIVE -> live
//	FT, AET, PEN, AWD, WO         -> finished
//	PST, SUSP, CANC, ABD          -> postponed
//
// If the short code is not recognised, ok is false and status is "".
// The caller decides how to handle unknown codes; this function never guesses.
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

// ----------------------------------------------------------------------------
// Domain conversion
// ----------------------------------------------------------------------------

// toDomain converts a raw API response item to a provider.Fixture.
func (it fixtureItem) toDomain() (provider.Fixture, error) {
	kickOff, err := time.Parse(time.RFC3339, it.Fixture.Date)
	if err != nil {
		return provider.Fixture{}, fmt.Errorf(
			"fixture %d: bad date %q: %w",
			it.Fixture.ID, it.Fixture.Date, err,
		)
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

	// Elapsed minutes are only meaningful when the match is live.
	if status == provider.StatusLive {
		f.Minute = it.Fixture.Status.Elapsed
	}

	return f, nil
}

// toDomainFixtures converts a slice of API items to domain fixtures.
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

// ----------------------------------------------------------------------------
// Provider interface implementation
// ----------------------------------------------------------------------------

// Leagues returns competitions for the given country.
// Pass an empty string to fetch all countries.
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

// Fixtures returns matches that match the given query.
// Date filters (From, To) are optional; a zero time means no filter.
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

// LiveFixtures returns matches that are currently in progress
// for the given league external IDs.
// Returns nil, nil when leagueExternalIDs is empty.
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
