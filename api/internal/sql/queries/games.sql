-- name: CreateGame :exec
INSERT INTO games (id, base_time, increment, white_id, black_id, rating_w, rating_b)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: InsertMove :exec
INSERT INTO moves (game_id, move_number, move_notation,orig, dest, move_fen, time_left)
VALUES ($1,$2, $3, $4, $5, $6, $7);

-- name: GetGameByID :one
SELECT id FROM games WHERE id = $1;

-- name: GetGameInfo :one
SELECT games.*, u1.username as white_username, u2.username as black_username FROM games
JOIN users u1 ON games.white_id = u1.id
JOIN users u2 ON games.black_id = u2.id
WHERE games.id = $1;

-- name: GetOngoingGames :many
SELECT * FROM games WHERE result = 'ongoing';

-- name: GetGameMoves :many
SELECT move_number, move_notation, orig, dest, move_fen, time_left
FROM moves
WHERE game_id = $1
ORDER BY move_number;

-- name: EndGameWithResult :exec
UPDATE games
SET result = $1, result_reason = $2, change_w = $3, change_b = $4, game_length = $5, end_time_left_white = $6, end_time_left_black = $7
WHERE id = $8;

-- name: GetPlayerGames :many
SELECT games.id, games.base_time, games.increment, u1.username as white_username, u2.username as black_username, games.result, games.game_length, games.result_reason, games.created_at, games.rating_w, games.rating_b, games.change_w, games.change_b
FROM games
JOIN users u1 ON games.white_id = u1.id
JOIN users u2 ON games.black_id = u2.id
WHERE (u1.username = $1 OR u2.username = $1)
ORDER BY games.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetGameNumbers :one
SELECT
COUNT(*) FILTER(WHERE games.white_id = users.id OR games.black_id = users.id) AS game_count,
COUNT(*) FILTER(WHERE (games.white_id = users.id AND result = '1-0') OR (games.black_id = users.id AND result = '0-1')) AS win_count,
COUNT(*) FILTER(WHERE (games.white_id = users.id OR games.black_id = users.id) AND result = '1/2-1/2') AS draw_count,
COUNT(*) FILTER(WHERE (games.white_id = users.id AND result = '0-1') OR (games.black_id = users.id AND result = '1-0')) AS loss_count
FROM games
JOIN users ON users.id = games.white_id or users.id = games.black_id
WHERE users.username = $1;

-- name: DeleteOngoingGames :exec
DELETE FROM games WHERE result = 'ongoing';