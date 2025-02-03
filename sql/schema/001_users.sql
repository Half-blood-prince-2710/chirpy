-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    email TEXT NOT NULL UNIQUE
);


-- +goose Down
DROP TABLE IF EXISTS users;