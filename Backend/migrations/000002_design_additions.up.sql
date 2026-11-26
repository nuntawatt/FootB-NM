-- แถบซ้ายของดีไซน์: แยกลีกกับถ้วย และกำหนดลำดับการแสดง
ALTER TABLE leagues
    ADD COLUMN type       VARCHAR(10) NOT NULL DEFAULT 'league' CHECK (type IN ('league', 'cup')),
    ADD COLUMN sort_order INT         NOT NULL DEFAULT 0;

-- "1-1 Pen" ในดีไซน์: คะแนนดวลจุดโทษ, นาทีที่แข่ง, รอบของถ้วย
ALTER TABLE matches
    ADD COLUMN home_penalty_score SMALLINT CHECK (home_penalty_score >= 0),
    ADD COLUMN away_penalty_score SMALLINT CHECK (away_penalty_score >= 0),
    ADD COLUMN minute             SMALLINT CHECK (minute >= 0),
    ADD COLUMN round              VARCHAR(50) NOT NULL DEFAULT '';

-- เมนู "ทีมโปรด"
CREATE TABLE user_favorite_teams (
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    team_id    UUID        NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, team_id)
);
CREATE INDEX idx_user_favorite_teams_team_id ON user_favorite_teams (team_id);
