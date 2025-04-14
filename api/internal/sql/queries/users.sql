-- name: CreateUser :one
INSERT INTO users(updated_at, email, avatar_url, google_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByUserID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUsernameByUserID :one
SELECT username FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUserAvatar :exec
UPDATE users SET avatar_url = $1, updated_at = $2 WHERE id = $3;

-- name: UpdateUsername :exec
UPDATE users SET username = $1 WHERE id = $2;

-- name: GetUserPublicInfo :one
SELECT created_at, avatar_url, rating, rd FROM users WHERE username = $1;

-- name: GetRatingInfo :one
SELECT rating, rd, volatility FROM users WHERE id = $1;

-- name: GetUsernameAndRating :one
SELECT username, rating FROM users WHERE id = $1;

-- name: UpdateRating :exec
UPDATE users SET rating = $1, rd = $2, volatility = $3 WHERE id = $4;