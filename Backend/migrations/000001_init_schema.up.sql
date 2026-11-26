CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(50)  NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password_hash TEXT         NOT NULL,
    avatar_url    TEXT         NOT NULL DEFAULT '',
    role          VARCHAR(10)  NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);
-- unique เฉพาะแถวที่ยังไม่ถูกลบ เพื่อให้สมัครใหม่ด้วยอีเมลเดิมได้หลัง soft delete
CREATE UNIQUE INDEX uq_users_username ON users (username)     WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX uq_users_email    ON users (lower(email)) WHERE deleted_at IS NULL;

CREATE TABLE leagues (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    short_name      VARCHAR(20)  NOT NULL DEFAULT '',
    country         VARCHAR(100) NOT NULL DEFAULT '',
    logo_url        TEXT         NOT NULL DEFAULT '',
    external_api_id VARCHAR(50),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE UNIQUE INDEX uq_leagues_external_api_id ON leagues (external_api_id)
    WHERE external_api_id IS NOT NULL AND deleted_at IS NULL;

CREATE TABLE teams (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    league_id       UUID         NOT NULL REFERENCES leagues (id) ON DELETE RESTRICT,
    name            VARCHAR(100) NOT NULL,
    short_name      VARCHAR(20)  NOT NULL DEFAULT '',
    logo_url        TEXT         NOT NULL DEFAULT '',
    external_api_id VARCHAR(50),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_teams_league_id ON teams (league_id);
CREATE UNIQUE INDEX uq_teams_external_api_id ON teams (external_api_id)
    WHERE external_api_id IS NOT NULL AND deleted_at IS NULL;

CREATE TABLE matches (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    league_id       UUID        NOT NULL REFERENCES leagues (id) ON DELETE RESTRICT,
    home_team_id    UUID        NOT NULL REFERENCES teams (id)   ON DELETE RESTRICT,
    away_team_id    UUID        NOT NULL REFERENCES teams (id)   ON DELETE RESTRICT,
    match_date      TIMESTAMPTZ NOT NULL,
    home_score      SMALLINT CHECK (home_score >= 0),
    away_score      SMALLINT CHECK (away_score >= 0),
    status          VARCHAR(20) NOT NULL DEFAULT 'scheduled'
                    CHECK (status IN ('scheduled', 'live', 'finished', 'postponed')),
    external_api_id VARCHAR(50),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CHECK (home_team_id <> away_team_id)
);
CREATE INDEX idx_matches_league_id  ON matches (league_id);
CREATE INDEX idx_matches_match_date ON matches (match_date);
CREATE INDEX idx_matches_status     ON matches (status);
CREATE UNIQUE INDEX uq_matches_external_api_id ON matches (external_api_id)
    WHERE external_api_id IS NOT NULL AND deleted_at IS NULL;

CREATE TABLE predictions (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id              UUID        NOT NULL REFERENCES users (id)   ON DELETE RESTRICT,
    match_id             UUID        NOT NULL REFERENCES matches (id) ON DELETE RESTRICT,
    predicted_home_score SMALLINT    NOT NULL CHECK (predicted_home_score >= 0),
    predicted_away_score SMALLINT    NOT NULL CHECK (predicted_away_score >= 0),
    predicted_result     VARCHAR(10) NOT NULL CHECK (predicted_result IN ('home_win', 'draw', 'away_win')),
    is_correct           BOOLEAN,
    points_earned        INT         NOT NULL DEFAULT 0,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at           TIMESTAMPTZ
);
CREATE UNIQUE INDEX uq_predictions_user_match ON predictions (user_id, match_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_predictions_match_id ON predictions (match_id);