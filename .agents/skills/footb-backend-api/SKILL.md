---
name: footb-backend-api
description: >-
  Cheatsheet of all backend REST API endpoints, data models, and JSON response
  shapes for the footb project. Use when building, debugging, or modifying
  the frontend that consumes these APIs.
---

# Football Backend API Reference

Base URL: `http://localhost:3000/api/v1`

## Endpoints (24 total)

| Category | Method | Path | Auth | Description / Params |
|----------|--------|------|------|----------------------|
| **Auth** | POST | `/register` (or `/users/register`) | Public | Body: `{ username, email, password, avatar_url? }` -> `{ status, token, user }` |
| | POST | `/login` (or `/users/login`) | Public | Body: `{ email, password }` (email or username) -> `{ status, token, user }` |
| | GET | `/profile` (or `/users/profile`) | Bearer | Get current user profile |
| | PUT | `/profile` (or `/users/profile`) | Bearer | Body: `{ username?, avatar_url? }` -> `{ status, message, user }` |
| **Seasons** | GET | `/seasons` | Public | List seasons 2022–2026 (`season`, `label`, `matches_count`) |
| **Leagues** | GET | `/leagues` | Public | List all leagues/cups |
| | GET | `/leagues/:id` | Public | Get league by UUID |
| **Teams** | GET | `/teams` | Public | List teams (`?league_id=UUID`) |
| | GET | `/teams/:id` | Public | Get team by UUID |
| | GET | `/teams/:id/form` | Public | Last 5 finished matches (`FormMatch[]`) |
| | GET | `/teams/:id/coach` | Public | Team head coach (`Coach`) |
| | GET | `/teams/:id/squad` | Public | Team players (`Player[]`) |
| **Matches** | GET | `/matches` | Public | List matches (`season`, `year`, `status`, `league_id`, `date`, `limit`, `page`) |
| | GET | `/matches/:id` | Public | Get match by UUID (preloaded) |
| | GET | `/matches/:id/lineups` | Public | Starting XI, formation, coach, subs (`MatchLineup`) |
| | GET | `/matches/:id/venue` | Public | Stadium, city, capacity, broadcast (`VenueInfo`) |
| **Standings** | GET | `/standings` | Public | Computed league table (`?season=2024&league_id=UUID`) |
| **Favorites** | GET | `/favorites` | Optional | List user's favorite `Team[]` (JWT or `X-User-ID` header) |
| | POST | `/favorites/:team_id` | Optional | Add team to favorites |
| | DELETE | `/favorites/:team_id` | Optional | Remove team from favorites |
| **Highlights** | GET | `/heroes` | Public | Homepage hero banners (`HeroHighlight[]`) |
| **System** | GET | `/health` | Public | Health check (`{ status: "ok" }`) |
| **Sync** | POST | `/sync/fixtures` | Public | Body: `{ season, from?, to? }` |
| | POST | `/sync/squads` | Public | Sync coaches & players for teams |

## Key Data Models (JSON shapes)

### User & Auth Response
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "avatar_url": "string",
  "role": "user | admin",
  "created_at": "timestamp"
}
```
Validation: `username` 3–30 (`[a-zA-Z0-9_]`), `password` 6–72 (must contain letter + number), `avatar_url` must start with `http://` or `https://` if set.

### League / Team / Match
- **League**: `id`, `name`, `short_name`, `country`, `logo_url`, `type` (`league|cup`), `sort_order`
- **Team**: `id`, `league_id`, `name`, `short_name`, `logo_url`, `league?`
- **Match**: `id`, `league_id`, `home_team_id`, `away_team_id`, `match_date`, `home_score?`, `away_score?`, `home_penalty_score?`, `away_penalty_score?`, `minute?`, `round`, `status` (`scheduled|live|finished|postponed`), `home_team`, `away_team`, `league`

### StandingsRow (`GET /standings`)
```json
{
  "rank": 1,
  "team_id": "uuid",
  "team_name": "string",
  "team_logo": "string",
  "played": 30,
  "won": 20,
  "drawn": 9,
  "lost": 1,
  "goals_for": 70,
  "goals_against": 27,
  "goal_diff": 43,
  "points": 69,
  "form": ["W", "W", "D", "W", "W"]
}
```

### FormMatch (`GET /teams/:id/form`)
```json
{
  "match_id": "uuid",
  "date": "YYYY-MM-DD",
  "opponent": "string",
  "score": "2 - 1",
  "result": "W | D | L",
  "is_home": true
}
```

### Coach (`GET /teams/:id/coach`) & Player (`GET /teams/:id/squad`)
- **Coach**: `id`, `team_id`, `name`, `nationality`, `photo_url`
- **Player**: `id`, `team_id`, `name`, `number?`, `position`, `age?`, `nationality`, `photo_url`

### MatchLineup (`GET /matches/:id/lineups`)
```json
{
  "match_id": "uuid",
  "home": {
    "team_id": "uuid",
    "team_name": "string",
    "formation": "4-3-3",
    "coach_name": "string",
    "starting": [{ "number": 1, "name": "string", "position": "GK", "rating": "7.0" }],
    "substitutes": []
  },
  "away": { "...": "same as home" }
}
```

### VenueInfo (`GET /matches/:id/venue`)
`{ "match_id": "uuid", "stadium": "string", "city": "string", "capacity": 25000, "broadcast": ["string"] }`

### HeroHighlight (`GET /heroes`)
`{ "id": "string", "title": "string", "subtitle": "string", "image_url": "string", "tag": "string", "match_id?": "uuid", "redirect_url?": "string" }`

## Key Notes for Frontend
1. All IDs are UUID strings; error responses return `{ "error": "message" }`.
2. Auth uses `Authorization: Bearer <JWT>` (7-day expiry).
3. `/favorites` supports JWT or `X-User-ID: <UUID>` fallback (defaults to `00000000-0000-0000-0000-000000000001` if unauthenticated).
4. `MatchLineup.starting[].rating` is a string (e.g. `"7.0"`), whereas `Team.squad` uses `Player` without `rating`.
