-- name: CreateGame :exec
INSERT INTO games (
                   id,
                   base_time,
                   increment,
                   tournament_id,
                   white_id,
                   black_id,
                   rating_w,
                   rating_b
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: InsertMoves :copyfrom
INSERT INTO moves (
                   game_id,
                   move_number,
                   move_notation,
                   orig,
                   dest,
                   move_fen,
                   time_left
)
VALUES ($1,$2, $3, $4, $5, $6, $7);

-- name: GetGameByID :one
SELECT id FROM games WHERE id = $1;

-- name: GetGameInfo :one
SELECT
    games.*,
    u1.username as white_username,
    u2.username as black_username,
    t.name as tournament_name
FROM games
LEFT OUTER JOIN tournaments t ON games.tournament_id = t.id
LEFT OUTER JOIN users u1 ON games.white_id = u1.id
LEFT OUTER JOIN users u2 ON games.black_id = u2.id
WHERE games.id = $1;

-- name: GetLiveGames :many
SELECT * FROM games WHERE result = 0;

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
SELECT t.name as tournament_name, t.id as tournament_id, games.id, games.base_time, games.increment, u1.username as white_username, u2.username as black_username, games.result, games.game_length, games.result_reason, games.created_at, games.rating_w, games.rating_b, games.change_w, games.change_b
FROM games
JOIN users u1 ON games.white_id = u1.id
JOIN users u2 ON games.black_id = u2.id
LEFT OUTER JOIN tournaments t ON games.tournament_id = t.id
WHERE (u1.username = $1 OR u2.username = $1)
ORDER BY games.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetGameNumbers :one
SELECT
COUNT(*) FILTER(WHERE games.white_id = users.id OR games.black_id = users.id) AS game_count,
COUNT(*) FILTER(WHERE (games.white_id = users.id AND result = 1) OR (games.black_id = users.id AND result = 2)) AS win_count,
COUNT(*) FILTER(WHERE (games.white_id = users.id OR games.black_id = users.id) AND result = 3) AS draw_count,
COUNT(*) FILTER(WHERE (games.white_id = users.id AND result = 2) OR (games.black_id = users.id AND result = 1)) AS loss_count
FROM games
JOIN users ON users.id = games.white_id or users.id = games.black_id
WHERE users.username = $1;

-- name: DeleteLiveGames :exec
DELETE FROM games WHERE result = 0;

-- name: CreateTournament :exec
INSERT INTO tournaments (id, name, start_time, duration, base_time, increment, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetTournamentByID :one
SELECT id FROM tournaments WHERE id = $1;

-- name: GetTournamentInfo :one
SELECT t.*, u.username FROM tournaments t JOIN users u ON t.created_by = u.id WHERE t.id = $1;

-- name: GetTournamentPlayers :many
SELECT tp.score, tp.streak, u.id, u.username, u.rating, tp.scores FROM tournament_players tp JOIN users u ON tp.player_id = u.id WHERE tp.tournament_id = $1;

-- name: InsertTournamentPlayer :exec
INSERT INTO tournament_players (player_id, tournament_id) VALUES ($1, $2);

-- name: UpdateTournamentStartTime :exec
UPDATE tournaments SET start_time = CURRENT_TIMESTAMP where id = $1;

-- name: GetUpcomingTournaments :many
SELECT * FROM tournaments WHERE status = 0 ORDER BY start_time;

-- name: GetTournamentPlayer :one
SELECT id FROM tournament_players WHERE player_id = $1 AND tournament_id = $2;

-- name: DeleteTournamentPlayer :exec
DELETE FROM tournament_players WHERE player_id = $1;

-- name: GetTournamentStatus :one
SELECT status FROM tournaments WHERE id = $1;

-- name: UpdateTournamentStatus :exec
UPDATE tournaments SET status = $1 WHERE id=$2;

-- name: BatchUpdateTournamentPlayers :exec
WITH players_data (
    id,
    score,
    scores,
    streak
    ) AS (
    SELECT
        (data->>'id')::int,
        (data->>'score')::int,
        (SELECT array_agg(value::smallint)
         FROM jsonb_array_elements(data->'scores') AS score_elements(value)),
        (data->>'streak')::int
    FROM jsonb_array_elements(@players_input::jsonb) AS data
        )
UPDATE tournament_players t
SET score = players_data.score,
    scores = players_data.scores,
    streak = players_data.streak
FROM players_data
WHERE t.tournament_id = $1 AND t.player_id = players_data.id;

-- name: GetTopN :many
SELECT id, username, avatar_url, rating FROM users
ORDER BY rating DESC LIMIT $1 OFFSET 0;