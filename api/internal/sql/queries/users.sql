-- name: CreateUser :one
INSERT INTO users(updated_at, email, avatar_url, google_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByUserID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUserAvatar :one
UPDATE users SET avatar_url = $1 WHERE id = $2 RETURNING *;