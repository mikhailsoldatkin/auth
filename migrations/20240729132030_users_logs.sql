-- +goose Up
CREATE TABLE users_logs
(
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT REFERENCES users (id) ON DELETE SET NULL,
    details   TEXT   NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users_logs;
