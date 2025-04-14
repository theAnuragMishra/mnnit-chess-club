-- name: CreateGame :one
INSERT INTO games (base_time, increment, white_id, black_id, white_username, black_username, fen, rating_w, rating_b)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: InsertMove :one
INSERT INTO moves (game_id, move_number, player_id, move_notation,orig, dest, move_fen)
VALUES ($1,$2, $3, $4, $5, $6, $7)
RETURNING move_number, move_notation, orig, dest, move_fen
;

-- name: GetGameInfo :one
SELECT * FROM games WHERE id = $1;

-- name: GetOngoingGames :many
SELECT * FROM games WHERE result = 'ongoing';

-- name: GetGameMoves :many
SELECT move_number, move_notation, orig, dest, move_fen
FROM moves
WHERE game_id = $1
ORDER BY move_number;

-- name: UpdateGameLengthAndFEN :exec
UPDATE games
SET fen = $1, game_length = $2
WHERE id = $3;

-- name: EndGameWithResult :exec
UPDATE games
SET result = $1, end_time_left_white = $2, end_time_left_black = $3, result_reason = $4, change_w = $5, change_b = $6
WHERE id = $7;

-- name: GetPlayerGames :many
SELECT id, base_time, increment, white_username, black_username, result, game_length, result_reason, created_at, rating_w, rating_b, change_w, change_b
FROM games
WHERE (white_username = $1 OR black_username = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetLatestMove :one
SELECT move_number, move_notation, orig, dest, move_fen
FROM moves
WHERE game_id = $1
ORDER BY move_number DESC
LIMIT $1;
