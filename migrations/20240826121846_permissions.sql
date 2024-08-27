-- +goose Up
CREATE TABLE permissions
(
    id       BIGSERIAL PRIMARY KEY,
    endpoint TEXT NOT NULL,
    role     TEXT NOT NULL,
    UNIQUE (endpoint, role)
);

-- predefined permissions
INSERT INTO permissions (endpoint, role)
VALUES ('/chat_v1.ChatV1/Create', 'ADMIN'),
       ('/chat_v1.ChatV1/Delete', 'ADMIN'),
       ('/chat_v1.ChatV1/SendMessage', 'ADMIN'),
       ('/chat_v1.ChatV1/SendMessage', 'USER');

-- +goose Down
DROP TABLE IF EXISTS permissions;
