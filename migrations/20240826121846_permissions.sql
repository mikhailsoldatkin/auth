-- +goose Up
CREATE TABLE permissions
(
    id       BIGSERIAL PRIMARY KEY,
    endpoint TEXT NOT NULL,
    role     TEXT NOT NULL,
    UNIQUE (endpoint, role)
);

-- Insert predefined permissions
INSERT INTO permissions (endpoint, role)
VALUES ('/user_v1.UserV1/Get', 'ADMIN'),
       ('/user_v1.UserV1/Get', 'USER'),
       ('/user_v1.UserV1/Create', 'ADMIN'),
       ('/user_v1.UserV1/Delete', 'ADMIN'),
       ('/user_v1.UserV1/List', 'ADMIN'),
       ('/user_v1.UserV1/List', 'USER');

-- +goose Down
DROP TABLE IF EXISTS permissions;
