-- name: CreateRefreshToken :one 
INSERT INTO refresh_tokens(token,user_id,expires_at)
VALUES ($1,$2,$3) RETURNING *;

