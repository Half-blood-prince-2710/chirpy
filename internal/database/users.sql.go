// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(email,hashed_password)
VALUES ($1,$2) RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, created_at,updated_at,email,hashed_password,is_chirpy_red FROM users WHERE email = $1
`

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateChirpRed = `-- name: UpdateChirpRed :one
UPDATE users
SET is_chirpy_red = $1 , updated_at = $2
WHERE id = $3 RETURNING id, created_at,updated_at,email, is_chirpy_red
`

type UpdateChirpRedParams struct {
	IsChirpyRed bool
	UpdatedAt   time.Time
	ID          uuid.UUID
}

type UpdateChirpRedRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	IsChirpyRed bool
}

func (q *Queries) UpdateChirpRed(ctx context.Context, arg UpdateChirpRedParams) (UpdateChirpRedRow, error) {
	row := q.db.QueryRowContext(ctx, updateChirpRed, arg.IsChirpyRed, arg.UpdatedAt, arg.ID)
	var i UpdateChirpRedRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users 
SET email = $1 , hashed_password = $2, updated_at = $3
where id = $4 RETURNING id, created_at,updated_at,email, is_chirpy_red
`

type UpdateUserParams struct {
	Email          string
	HashedPassword string
	UpdatedAt      time.Time
	ID             uuid.UUID
}

type UpdateUserRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	IsChirpyRed bool
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Email,
		arg.HashedPassword,
		arg.UpdatedAt,
		arg.ID,
	)
	var i UpdateUserRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}
