// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: games.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createGame = `-- name: CreateGame :one
INSERT INTO games (white_player_id, black_player_id)
VALUES ($1, $2)
    RETURNING id
`

type CreateGameParams struct {
	WhitePlayerID pgtype.UUID
	BlackPlayerID pgtype.UUID
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) (int32, error) {
	row := q.db.QueryRow(ctx, createGame, arg.WhitePlayerID, arg.BlackPlayerID)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const endGameWithResult = `-- name: EndGameWithResult :exec
UPDATE games
SET result = '1-0', ended_at = NOW()
WHERE id = 1
`

func (q *Queries) EndGameWithResult(ctx context.Context) error {
	_, err := q.db.Exec(ctx, endGameWithResult)
	return err
}

const getGameMoves = `-- name: GetGameMoves :many
SELECT move_number, move_notation, move_fen
FROM moves
WHERE game_id = $1
ORDER BY move_number
`

type GetGameMovesRow struct {
	MoveNumber   int32
	MoveNotation string
	MoveFen      string
}

func (q *Queries) GetGameMoves(ctx context.Context, gameID pgtype.Int4) ([]GetGameMovesRow, error) {
	rows, err := q.db.Query(ctx, getGameMoves, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGameMovesRow
	for rows.Next() {
		var i GetGameMovesRow
		if err := rows.Scan(&i.MoveNumber, &i.MoveNotation, &i.MoveFen); err != nil {
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
SELECT move_number, move_notation, move_fen
FROM moves
WHERE game_id = 1
ORDER BY move_number DESC
LIMIT 1
`

type GetLatestMoveRow struct {
	MoveNumber   int32
	MoveNotation string
	MoveFen      string
}

func (q *Queries) GetLatestMove(ctx context.Context) (GetLatestMoveRow, error) {
	row := q.db.QueryRow(ctx, getLatestMove)
	var i GetLatestMoveRow
	err := row.Scan(&i.MoveNumber, &i.MoveNotation, &i.MoveFen)
	return i, err
}

const getPlayerGames = `-- name: GetPlayerGames :many
SELECT id, white_player_id, black_player_id, result, ended_at
FROM games
WHERE (white_player_id = 1 OR black_player_id = 1)
AND ended_at IS NOT NULL
ORDER BY ended_at DESC
`

type GetPlayerGamesRow struct {
	ID            int32
	WhitePlayerID pgtype.UUID
	BlackPlayerID pgtype.UUID
	Result        pgtype.Text
	EndedAt       pgtype.Timestamp
}

func (q *Queries) GetPlayerGames(ctx context.Context) ([]GetPlayerGamesRow, error) {
	rows, err := q.db.Query(ctx, getPlayerGames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlayerGamesRow
	for rows.Next() {
		var i GetPlayerGamesRow
		if err := rows.Scan(
			&i.ID,
			&i.WhitePlayerID,
			&i.BlackPlayerID,
			&i.Result,
			&i.EndedAt,
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

const insertMove = `-- name: InsertMove :exec
INSERT INTO moves (game_id, move_number, player_id, move_notation, move_fen)
VALUES ($1,$2, $3, $4, $5)
`

type InsertMoveParams struct {
	GameID       pgtype.Int4
	MoveNumber   int32
	PlayerID     pgtype.UUID
	MoveNotation string
	MoveFen      string
}

func (q *Queries) InsertMove(ctx context.Context, arg InsertMoveParams) error {
	_, err := q.db.Exec(ctx, insertMove,
		arg.GameID,
		arg.MoveNumber,
		arg.PlayerID,
		arg.MoveNotation,
		arg.MoveFen,
	)
	return err
}
