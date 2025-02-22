-- name: CreateGame :one
INSERT INTO games (base_time, increment, white_id, black_id, white_username, black_username, fen)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: InsertMove :one
INSERT INTO moves (game_id, move_number, player_id, move_notation, move_fen)
VALUES ($1,$2, $3, $4, $5)
RETURNING move_number, move_notation, move_fen
;

-- name: GetGameInfo :one
SELECT * FROM games WHERE id = $1;

-- name: GetOngoingGames :many
SELECT * FROM games WHERE result = 'ongoing';

-- name: GetGameMoves :many
SELECT move_number, move_notation, move_fen
FROM moves
WHERE game_id = $1
ORDER BY move_number;

-- name: UpdateGameLengthAndFEN :exec
UPDATE games
SET fen = $1, game_length = $2
WHERE id = $3;

-- name: EndGameWithResult :exec
UPDATE games
SET result = $1
WHERE id = $2;

-- name: GetPlayerGames :many
SELECT id, white_username, black_username, result
FROM games
WHERE (white_username = $1 OR black_username = $1)
--AND ended_at IS NOT NULL
ORDER BY created_at DESC;

-- name: GetLatestMove :one
SELECT move_number, move_notation, move_fen
FROM moves
WHERE game_id = $1
ORDER BY move_number DESC
LIMIT $1;
