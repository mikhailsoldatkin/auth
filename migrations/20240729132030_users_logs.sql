-- +goose Up
CREATE TABLE users_logs
(
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    details   TEXT   NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users_logs;
