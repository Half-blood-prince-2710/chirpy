-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_tokens (
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    expires_at TIMESTAMP,
    revoked_at TIMESTAMP,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);




-- +goose Down 
DROP TABLE IF EXISTS refresh_tokens;