// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: games.sql

package database

import (
	"context"
	"time"
)

const createGame = `-- name: CreateGame :exec
INSERT INTO games (id, base_time, increment, white_id, black_id, rating_w, rating_b)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type CreateGameParams struct {
	ID        string
	BaseTime  int32
	Increment int32
	WhiteID   *int32
	BlackID   *int32
	RatingW   int32
	RatingB   int32
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) error {
	_, err := q.db.Exec(ctx, createGame,
		arg.ID,
		arg.BaseTime,
		arg.Increment,
		arg.WhiteID,
		arg.BlackID,
		arg.RatingW,
		arg.RatingB,
	)
	return err
}

const deleteOngoingGames = `-- name: DeleteOngoingGames :exec
DELETE FROM games WHERE result = 'ongoing'
`

func (q *Queries) DeleteOngoingGames(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deleteOngoingGames)
	return err
}

const endGameWithResult = `-- name: EndGameWithResult :exec
UPDATE games
SET result = $1, result_reason = $2, change_w = $3, change_b = $4, game_length = $5, end_time_left_white = $6, end_time_left_black = $7
WHERE id = $8
`

type EndGameWithResultParams struct {
	Result           string
	ResultReason     *string
	ChangeW          *int32
	ChangeB          *int32
	GameLength       int16
	EndTimeLeftWhite *int32
	EndTimeLeftBlack *int32
	ID               string
}

func (q *Queries) EndGameWithResult(ctx context.Context, arg EndGameWithResultParams) error {
	_, err := q.db.Exec(ctx, endGameWithResult,
		arg.Result,
		arg.ResultReason,
		arg.ChangeW,
		arg.ChangeB,
		arg.GameLength,
		arg.EndTimeLeftWhite,
		arg.EndTimeLeftBlack,
		arg.ID,
	)
	return err
}

const getGameByID = `-- name: GetGameByID :one
SELECT id FROM games WHERE id = $1
`

func (q *Queries) GetGameByID(ctx context.Context, id string) (string, error) {
	row := q.db.QueryRow(ctx, getGameByID, id)
	err := row.Scan(&id)
	return id, err
}

const getGameInfo = `-- name: GetGameInfo :one
SELECT games.id, games.base_time, games.increment, games.white_id, games.black_id, games.game_length, games.result, games.created_at, games.end_time_left_white, games.end_time_left_black, games.result_reason, games.rating_w, games.rating_b, games.change_w, games.change_b, u1.username as white_username, u2.username as black_username FROM games
JOIN users u1 ON games.white_id = u1.id
JOIN users u2 ON games.black_id = u2.id
WHERE games.id = $1
`

type GetGameInfoRow struct {
	ID               string
	BaseTime         int32
	Increment        int32
	WhiteID          *int32
	BlackID          *int32
	GameLength       int16
	Result           string
	CreatedAt        time.Time
	EndTimeLeftWhite *int32
	EndTimeLeftBlack *int32
	ResultReason     *string
	RatingW          int32
	RatingB          int32
	ChangeW          *int32
	ChangeB          *int32
	WhiteUsername    *string
	BlackUsername    *string
}

func (q *Queries) GetGameInfo(ctx context.Context, id string) (GetGameInfoRow, error) {
	row := q.db.QueryRow(ctx, getGameInfo, id)
	var i GetGameInfoRow
	err := row.Scan(
		&i.ID,
		&i.BaseTime,
		&i.Increment,
		&i.WhiteID,
		&i.BlackID,
		&i.GameLength,
		&i.Result,
		&i.CreatedAt,
		&i.EndTimeLeftWhite,
		&i.EndTimeLeftBlack,
		&i.ResultReason,
		&i.RatingW,
		&i.RatingB,
		&i.ChangeW,
		&i.ChangeB,
		&i.WhiteUsername,
		&i.BlackUsername,
	)
	return i, err
}

const getGameMoves = `-- name: GetGameMoves :many
SELECT move_number, move_notation, orig, dest, move_fen, time_left
FROM moves
WHERE game_id = $1
ORDER BY move_number
`

type GetGameMovesRow struct {
	MoveNumber   int32
	MoveNotation string
	Orig         string
	Dest         string
	MoveFen      string
	TimeLeft     *int32
}

func (q *Queries) GetGameMoves(ctx context.Context, gameID string) ([]GetGameMovesRow, error) {
	rows, err := q.db.Query(ctx, getGameMoves, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGameMovesRow
	for rows.Next() {
		var i GetGameMovesRow
		if err := rows.Scan(
			&i.MoveNumber,
			&i.MoveNotation,
			&i.Orig,
			&i.Dest,
			&i.MoveFen,
			&i.TimeLeft,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGameNumbers = `-- name: GetGameNumbers :one
SELECT
COUNT(*) FILTER(WHERE games.white_id = users.id OR games.black_id = users.id) AS game_count,
COUNT(*) FILTER(WHERE (games.white_id = users.id AND result = '1-0') OR (games.black_id = users.id AND result = '0-1')) AS win_count,
COUNT(*) FILTER(WHERE (games.white_id = users.id OR games.black_id = users.id) AND result = '1/2-1/2') AS draw_count,
COUNT(*) FILTER(WHERE (games.white_id = users.id AND result = '0-1') OR (games.black_id = users.id AND result = '1-0')) AS loss_count
FROM games
JOIN users ON users.id = games.white_id or users.id = games.black_id
WHERE users.username = $1
`

type GetGameNumbersRow struct {
	GameCount int64
	WinCount  int64
	DrawCount int64
	LossCount int64
}

func (q *Queries) GetGameNumbers(ctx context.Context, username *string) (GetGameNumbersRow, error) {
	row := q.db.QueryRow(ctx, getGameNumbers, username)
	var i GetGameNumbersRow
	err := row.Scan(
		&i.GameCount,
		&i.WinCount,
		&i.DrawCount,
		&i.LossCount,
	)
	return i, err
}

const getOngoingGames = `-- name: GetOngoingGames :many
SELECT id, base_time, increment, white_id, black_id, game_length, result, created_at, end_time_left_white, end_time_left_black, result_reason, rating_w, rating_b, change_w, change_b FROM games WHERE result = 'ongoing'
`

func (q *Queries) GetOngoingGames(ctx context.Context) ([]Game, error) {
	rows, err := q.db.Query(ctx, getOngoingGames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Game
	for rows.Next() {
		var i Game
		if err := rows.Scan(
			&i.ID,
			&i.BaseTime,
			&i.Increment,
			&i.WhiteID,
			&i.BlackID,
			&i.GameLength,
			&i.Result,
			&i.CreatedAt,
			&i.EndTimeLeftWhite,
			&i.EndTimeLeftBlack,
			&i.ResultReason,
			&i.RatingW,
			&i.RatingB,
			&i.ChangeW,
			&i.ChangeB,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlayerGames = `-- name: GetPlayerGames :many
SELECT games.id, games.base_time, games.increment, u1.username as white_username, u2.username as black_username, games.result, games.game_length, games.result_reason, games.created_at, games.rating_w, games.rating_b, games.change_w, games.change_b
FROM games
JOIN users u1 ON games.white_id = u1.id
JOIN users u2 ON games.black_id = u2.id
WHERE (u1.username = $1 OR u2.username = $1)
ORDER BY games.created_at DESC
LIMIT $2 OFFSET $3
`

type GetPlayerGamesParams struct {
	Username *string
	Limit    int32
	Offset   int32
}

type GetPlayerGamesRow struct {
	ID            string
	BaseTime      int32
	Increment     int32
	WhiteUsername *string
	BlackUsername *string
	Result        string
	GameLength    int16
	ResultReason  *string
	CreatedAt     time.Time
	RatingW       int32
	RatingB       int32
	ChangeW       *int32
	ChangeB       *int32
}

func (q *Queries) GetPlayerGames(ctx context.Context, arg GetPlayerGamesParams) ([]GetPlayerGamesRow, error) {
	rows, err := q.db.Query(ctx, getPlayerGames, arg.Username, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlayerGamesRow
	for rows.Next() {
		var i GetPlayerGamesRow
		if err := rows.Scan(
			&i.ID,
			&i.BaseTime,
			&i.Increment,
			&i.WhiteUsername,
			&i.BlackUsername,
			&i.Result,
			&i.GameLength,
			&i.ResultReason,
			&i.CreatedAt,
			&i.RatingW,
			&i.RatingB,
			&i.ChangeW,
			&i.ChangeB,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type InsertMovesParams struct {
	GameID       string
	MoveNumber   int32
	MoveNotation string
	Orig         string
	Dest         string
	MoveFen      string
	TimeLeft     *int32
}
