// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: games.sql

package database

import (
	"context"
	"time"
)

const createGame = `-- name: CreateGame :one
INSERT INTO games (base_time, increment, white_id, black_id, white_username, black_username, fen, rating_w, rating_b)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id
`

type CreateGameParams struct {
	BaseTime      int32
	Increment     int32
	WhiteID       *int32
	BlackID       *int32
	WhiteUsername *string
	BlackUsername *string
	Fen           string
	RatingW       int32
	RatingB       int32
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) (int32, error) {
	row := q.db.QueryRow(ctx, createGame,
		arg.BaseTime,
		arg.Increment,
		arg.WhiteID,
		arg.BlackID,
		arg.WhiteUsername,
		arg.BlackUsername,
		arg.Fen,
		arg.RatingW,
		arg.RatingB,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const endGameWithResult = `-- name: EndGameWithResult :exec
UPDATE games
SET result = $1, end_time_left_white = $2, end_time_left_black = $3, result_reason = $4, change_w = $5, change_b = $6
WHERE id = $7
`

type EndGameWithResultParams struct {
	Result           string
	EndTimeLeftWhite *int32
	EndTimeLeftBlack *int32
	ResultReason     *string
	ChangeW          *int32
	ChangeB          *int32
	ID               int32
}

func (q *Queries) EndGameWithResult(ctx context.Context, arg EndGameWithResultParams) error {
	_, err := q.db.Exec(ctx, endGameWithResult,
		arg.Result,
		arg.EndTimeLeftWhite,
		arg.EndTimeLeftBlack,
		arg.ResultReason,
		arg.ChangeW,
		arg.ChangeB,
		arg.ID,
	)
	return err
}

const getGameInfo = `-- name: GetGameInfo :one
SELECT id, base_time, increment, white_id, black_id, white_username, black_username, fen, game_length, result, created_at, end_time_left_white, end_time_left_black, result_reason, rating_w, rating_b, change_w, change_b FROM games WHERE id = $1
`

func (q *Queries) GetGameInfo(ctx context.Context, id int32) (Game, error) {
	row := q.db.QueryRow(ctx, getGameInfo, id)
	var i Game
	err := row.Scan(
		&i.ID,
		&i.BaseTime,
		&i.Increment,
		&i.WhiteID,
		&i.BlackID,
		&i.WhiteUsername,
		&i.BlackUsername,
		&i.Fen,
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
	)
	return i, err
}

const getGameMoves = `-- name: GetGameMoves :many
SELECT move_number, move_notation, orig, dest, move_fen
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
}

func (q *Queries) GetGameMoves(ctx context.Context, gameID int32) ([]GetGameMovesRow, error) {
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

const getLatestMove = `-- name: GetLatestMove :one
SELECT move_number, move_notation, orig, dest, move_fen
FROM moves
WHERE game_id = $1
ORDER BY move_number DESC
LIMIT $1
`

type GetLatestMoveRow struct {
	MoveNumber   int32
	MoveNotation string
	Orig         string
	Dest         string
	MoveFen      string
}

func (q *Queries) GetLatestMove(ctx context.Context, limit int32) (GetLatestMoveRow, error) {
	row := q.db.QueryRow(ctx, getLatestMove, limit)
	var i GetLatestMoveRow
	err := row.Scan(
		&i.MoveNumber,
		&i.MoveNotation,
		&i.Orig,
		&i.Dest,
		&i.MoveFen,
	)
	return i, err
}

const getOngoingGames = `-- name: GetOngoingGames :many
SELECT id, base_time, increment, white_id, black_id, white_username, black_username, fen, game_length, result, created_at, end_time_left_white, end_time_left_black, result_reason, rating_w, rating_b, change_w, change_b FROM games WHERE result = 'ongoing'
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
			&i.WhiteUsername,
			&i.BlackUsername,
			&i.Fen,
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
SELECT id, base_time, increment, white_username, black_username, result, game_length, result_reason, created_at, rating_w, rating_b, change_w, change_b
FROM games
WHERE (white_username = $1 OR black_username = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetPlayerGamesParams struct {
	WhiteUsername *string
	Limit         int32
	Offset        int32
}

type GetPlayerGamesRow struct {
	ID            int32
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
	rows, err := q.db.Query(ctx, getPlayerGames, arg.WhiteUsername, arg.Limit, arg.Offset)
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

const insertMove = `-- name: InsertMove :one
INSERT INTO moves (game_id, move_number, player_id, move_notation,orig, dest, move_fen)
VALUES ($1,$2, $3, $4, $5, $6, $7)
RETURNING move_number, move_notation, orig, dest, move_fen
`

type InsertMoveParams struct {
	GameID       int32
	MoveNumber   int32
	PlayerID     *int32
	MoveNotation string
	Orig         string
	Dest         string
	MoveFen      string
}

type InsertMoveRow struct {
	MoveNumber   int32
	MoveNotation string
	Orig         string
	Dest         string
	MoveFen      string
}

func (q *Queries) InsertMove(ctx context.Context, arg InsertMoveParams) (InsertMoveRow, error) {
	row := q.db.QueryRow(ctx, insertMove,
		arg.GameID,
		arg.MoveNumber,
		arg.PlayerID,
		arg.MoveNotation,
		arg.Orig,
		arg.Dest,
		arg.MoveFen,
	)
	var i InsertMoveRow
	err := row.Scan(
		&i.MoveNumber,
		&i.MoveNotation,
		&i.Orig,
		&i.Dest,
		&i.MoveFen,
	)
	return i, err
}

const updateGameLengthAndFEN = `-- name: UpdateGameLengthAndFEN :exec
UPDATE games
SET fen = $1, game_length = $2
WHERE id = $3
`

type UpdateGameLengthAndFENParams struct {
	Fen        string
	GameLength int16
	ID         int32
}

func (q *Queries) UpdateGameLengthAndFEN(ctx context.Context, arg UpdateGameLengthAndFENParams) error {
	_, err := q.db.Exec(ctx, updateGameLengthAndFEN, arg.Fen, arg.GameLength, arg.ID)
	return err
}
