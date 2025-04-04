// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(updated_at, email, avatar_url, google_id)
VALUES ($1, $2, $3, $4)
RETURNING id, email, created_at, updated_at, username, avatar_url, google_id
`

type CreateUserParams struct {
	UpdatedAt time.Time
	Email     string
	AvatarUrl *string
	GoogleID  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.UpdatedAt,
		arg.Email,
		arg.AvatarUrl,
		arg.GoogleID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.AvatarUrl,
		&i.GoogleID,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, created_at, updated_at, username, avatar_url, google_id FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.AvatarUrl,
		&i.GoogleID,
	)
	return i, err
}

const getUserByUserID = `-- name: GetUserByUserID :one
SELECT id, email, created_at, updated_at, username, avatar_url, google_id FROM users WHERE id = $1
`

func (q *Queries) GetUserByUserID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUserID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.AvatarUrl,
		&i.GoogleID,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, email, created_at, updated_at, username, avatar_url, google_id FROM users WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username *string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.AvatarUrl,
		&i.GoogleID,
	)
	return i, err
}

const updateUserAvatar = `-- name: UpdateUserAvatar :one
UPDATE users SET avatar_url = $1 WHERE id = $2 RETURNING id, email, created_at, updated_at, username, avatar_url, google_id
`

type UpdateUserAvatarParams struct {
	AvatarUrl *string
	ID        int32
}

func (q *Queries) UpdateUserAvatar(ctx context.Context, arg UpdateUserAvatarParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserAvatar, arg.AvatarUrl, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.AvatarUrl,
		&i.GoogleID,
	)
	return i, err
}

const updateUsername = `-- name: UpdateUsername :exec
UPDATE users SET username = $1 WHERE id = $2
`

type UpdateUsernameParams struct {
	Username *string
	ID       int32
}

func (q *Queries) UpdateUsername(ctx context.Context, arg UpdateUsernameParams) error {
	_, err := q.db.Exec(ctx, updateUsername, arg.Username, arg.ID)
	return err
}
