-- name: CreateUser :one
INSERT INTO users(email,hashed_password)
VALUES ($1,$2) RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: FindUserByEmail :one
SELECT id, created_at,updated_at,email,hashed_password,is_chirpy_red FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users 
SET email = $1 , hashed_password = $2, updated_at = $3
where id = $4 RETURNING id, created_at,updated_at,email, is_chirpy_red;

-- name: UpdateChirpRed :one
UPDATE users
SET is_chirpy_red = $1 , updated_at = $2
WHERE id = $3 RETURNING id, created_at,updated_at,email, is_chirpy_red;