CREATE TABLE coaches (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID         NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    nationality     VARCHAR(100) NOT NULL DEFAULT '',
    photo_url       TEXT         NOT NULL DEFAULT '',
    external_api_id VARCHAR(50),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_coaches_team_id ON coaches (team_id);
CREATE UNIQUE INDEX uq_coaches_external_api_id ON coaches (external_api_id)
    WHERE external_api_id IS NOT NULL AND deleted_at IS NULL;

CREATE TABLE players (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID         NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    number          INT,
    position        VARCHAR(50)  NOT NULL DEFAULT '',
    age             INT,
    nationality     VARCHAR(100) NOT NULL DEFAULT '',
    photo_url       TEXT         NOT NULL DEFAULT '',
    external_api_id VARCHAR(50),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_players_team_id ON players (team_id);
CREATE UNIQUE INDEX uq_players_external_api_id ON players (external_api_id)
    WHERE external_api_id IS NOT NULL AND deleted_at IS NULL;
