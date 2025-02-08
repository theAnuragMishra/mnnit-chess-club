-- name: CreateGame :one
INSERT INTO games (white_player_id, black_player_id)
VALUES ($1, $2)
    RETURNING id;

-- name: InsertMove :exec
INSERT INTO moves (game_id, move_number, player_id, move_notation, move_fen)
VALUES ($1,$2, $3, $4, $5);

-- name: GetGameMoves :many
SELECT move_number, move_notation, move_fen
FROM moves
WHERE game_id = $1
ORDER BY move_number;

-- name: EndGameWithResult :exec
UPDATE games
SET result = '1-0', ended_at = NOW()
WHERE id = 1;

-- name: GetPlayerGames :many
SELECT id, white_player_id, black_player_id, result, ended_at
FROM games
WHERE (white_player_id = 1 OR black_player_id = 1)
AND ended_at IS NOT NULL
ORDER BY ended_at DESC;

-- name: GetLatestMove :one
SELECT move_number, move_notation, move_fen
FROM moves
WHERE game_id = 1
ORDER BY move_number DESC
LIMIT 1;
