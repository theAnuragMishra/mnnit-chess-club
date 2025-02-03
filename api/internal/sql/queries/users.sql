-- name: CreateUser :exec
INSERT INTO users(id, created_at, updated_at, username, password_hash)
VALUES ($1, $2, $3, $4, $5);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByUserID :one
SELECT * FROM users WHERE id = $1;