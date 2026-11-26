---
name: footb-backend-api
description: >-
  Cheatsheet of all backend REST API endpoints, data models, and JSON response
  shapes for the footb project. Use when building, debugging, or modifying
  the frontend that consumes these APIs.
---

# Football Backend API Reference

Base URL: `http://localhost:3000/api/v1`

## Endpoints

| # | Method | Path | Description | Query Params |
|---|--------|------|-------------|--------------|
| 1 | GET | `/health` | Health check (server + DB) | — |
| 2 | GET | `/seasons` | List seasons 2022–2026 for dropdown | — |
| 3 | GET | `/leagues` | List all leagues/cups (ordered by name) | — |
| 4 | GET | `/leagues/:id` | Get league by UUID | — |
| 5 | GET | `/teams` | List teams (ordered by name) | `league_id` (UUID) |
| 6 | GET | `/teams/:id` | Get team by UUID | — |
| 7 | GET | `/matches` | List matches with preloaded relations | `season`, `year`, `status`, `league_id`, `date` (YYYY-MM-DD), `limit` (1–200, default 50), `page` |
| 8 | GET | `/matches/:id` | Get match by UUID (preloaded) | — |
| 9 | POST | `/sync/fixtures` | Sync data from API-Football | Body: `{ "season": int, "from?": "YYYY-MM-DD", "to?": "YYYY-MM-DD" }` |

## Data Models (JSON shapes)

### League
```json
{
  "id": "uuid",
  "name": "string",
  "short_name": "string",
  "country": "string",
  "logo_url": "string",
  "type": "league | cup",
  "sort_order": 0,
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Team
```json
{
  "id": "uuid",
  "league_id": "uuid",
  "name": "string",
  "short_name": "string",
  "logo_url": "string",
  "league": "League? (omitempty)",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Match (preloaded)
```json
{
  "id": "uuid",
  "league_id": "uuid",
  "home_team_id": "uuid",
  "away_team_id": "uuid",
  "match_date": "timestamp",
  "home_score": "int16?",
  "away_score": "int16?",
  "home_penalty_score": "int16?",
  "away_penalty_score": "int16?",
  "minute": "int16?",
  "round": "string",
  "status": "scheduled | live | finished | postponed",
  "home_team": "Team",
  "away_team": "Team",
  "league": "League",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Season Item (from /seasons)
```json
{
  "season": 2024,
  "label": "2024/2025",
  "matches_count": 240
}
```

### User (future — no API endpoint yet)
Fields: `id`, `username`, `email`, `avatar_url`, `role` (user|admin)
`password_hash` hidden from JSON.

### Prediction (future — no API endpoint yet)
Fields: `id`, `user_id`, `match_id`, `predicted_home_score`, `predicted_away_score`,
`predicted_result` (home_win|draw|away_win), `is_correct`, `points_earned`

### UserFavoriteTeam (future — no API endpoint yet)
Fields: `user_id`, `team_id`, `created_at`

## Tech Stack
- **Language:** Go
- **Framework:** Fiber v3
- **ORM:** GORM (PostgreSQL)
- **IDs:** UUID (gen_random_uuid)
- **Soft delete:** `deleted_at` column (GORM convention)
- **CORS:** enabled globally (open)
- **Error format:** `{ "error": "message" }` with appropriate HTTP status code
- **Config:** ENV vars via `.env` — `PORT` (default 3000), `DATABASE_URL`, `FOOTBALL_API_KEY`

## Key Notes for Frontend
1. All IDs are UUID strings
2. Nullable numeric fields (`home_score`, `away_score`, etc.) return `null` when not set
3. Match list ordered by `match_date DESC`
4. Season filter: `season=2024` = Jul 2024 – Jun 2025
5. Year filter: `year=2024` = Jan 2024 – Dec 2024
6. Default page size 50, max 200
7. Preloaded relations (`home_team`, `away_team`, `league`) included in match responses
8. No auth required for any endpoint currently
