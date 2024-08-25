-- +goose Up
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    username   TEXT                     NOT NULL,
    email      TEXT UNIQUE              NOT NULL,
    role       TEXT                     NOT NULL DEFAULT 'UNKNOWN',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
