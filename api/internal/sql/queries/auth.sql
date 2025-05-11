-- name: CreateSession :exec
INSERT INTO sessions(id, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: GetSession :one
SELECT sessions.*,users.username FROM sessions JOIN users ON sessions.user_id = users.id WHERE sessions.id = $1;

-- name: UpdateSessionExpiry :exec
UPDATE sessions SET expires_at = $1 WHERE id = $2;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;
