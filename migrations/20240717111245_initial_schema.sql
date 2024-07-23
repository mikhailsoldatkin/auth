-- +goose Up
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT                     NOT NULL,
    email      TEXT UNIQUE              NOT NULL,
    role       TEXT                     NOT NULL DEFAULT 'UNKNOWN',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
