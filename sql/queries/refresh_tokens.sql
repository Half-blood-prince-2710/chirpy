-- name: CreateRefreshToken :one 
INSERT INTO refresh_tokens(token,user_id,expires_at)
VALUES ($1,$2,$3) RETURNING *;

-- name: GetTokenByToken :one
SELECT * FROM refresh_tokens WHERE token = $1;