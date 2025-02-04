-- name: CreateUser :one
INSERT INTO users(email,hashed_password)
VALUES ($1,$2) RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;