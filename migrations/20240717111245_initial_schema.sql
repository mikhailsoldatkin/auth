-- +goose Up
CREATE TYPE user_role AS ENUM ('UNKNOWN', 'USER', 'ADMIN');

CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255)             NOT NULL,
    email      VARCHAR(255) UNIQUE      NOT NULL,
    role       user_role                NOT NULL DEFAULT 'UNKNOWN',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
