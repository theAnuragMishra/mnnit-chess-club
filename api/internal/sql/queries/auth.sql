-- name: CreateSession :exec
INSERT INTO sessions(id, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: CreateCSRFToken :exec
INSERT INTO csrf_tokens(session_id, token, expires_at)
VALUES ($1, $2, $3);

-- name: GetSession :one
SELECT sessions.*,users.username FROM sessions JOIN users ON sessions.user_id = users.id WHERE sessions.id = $1;

-- name: UpdateSessionExpiry :exec
UPDATE sessions SET expires_at = $1 WHERE id = $2;

-- name: UpdateCSRFToken :exec
UPDATE csrf_tokens SET token = $1 WHERE session_id = $2;

-- name: GetCSRFTokenBySession :one
SELECT * from csrf_tokens WHERE session_id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;
