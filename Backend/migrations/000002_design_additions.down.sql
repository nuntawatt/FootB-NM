DROP TABLE IF EXISTS user_favorite_teams;

ALTER TABLE matches
    DROP COLUMN IF EXISTS round,
    DROP COLUMN IF EXISTS minute,
    DROP COLUMN IF EXISTS away_penalty_score,
    DROP COLUMN IF EXISTS home_penalty_score;

ALTER TABLE leagues
    DROP COLUMN IF EXISTS sort_order,
    DROP COLUMN IF EXISTS type;