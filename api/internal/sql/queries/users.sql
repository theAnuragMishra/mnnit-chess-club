-- name: CreateUser :exec
INSERT INTO users(updated_at, username, password_hash)
VALUES ($1, $2, $3);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByUserID :one
SELECT * FROM users WHERE id = $1;