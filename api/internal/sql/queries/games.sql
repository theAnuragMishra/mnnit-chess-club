-- name: CreateGame :one
INSERT INTO games (white_id, black_id, white_username, black_username)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: InsertMove :exec
INSERT INTO moves (game_id, move_number, player_id, move_notation, move_fen)
VALUES ($1,$2, $3, $4, $5);

-- name: GetGameInfo :one
SELECT * FROM games WHERE id = $1;

-- name: GetGameMoves :many
SELECT move_number, move_notation, move_fen
FROM moves
WHERE game_id = $1
ORDER BY move_number;

-- name: EndGameWithResult :exec
UPDATE games
SET result = $1, ended_at = NOW()
WHERE id = $2;

-- name: GetPlayerGames :many
SELECT id, white_username, black_username, result
FROM games
WHERE (white_username = $1 OR black_username = $1);
--AND ended_at IS NOT NULL
--ORDER BY ended_at DESC;

-- name: GetLatestMove :one
SELECT move_number, move_notation, move_fen
FROM moves
WHERE game_id = $1
ORDER BY move_number DESC
LIMIT $1;
