package game

import (
	"log"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"

	"github.com/notnil/chess"
)

type Game struct {
	ID            int32
	Result        string
	BaseTime      time.Duration
	Increment     time.Duration
	WhiteID       int32
	BlackID       int32
	Board         *chess.Game
	GameLength    int16
	LastMoveTime  time.Time
	TimeWhite     time.Duration
	TimeBlack     time.Duration
	DrawOfferedBy int32
	Timer         *time.Timer
}

func NewGame(baseTime time.Duration, increment time.Duration, player1 int32, player2 int32) *Game {
	board := chess.NewGame()

	return &Game{
		Result:       "ongoing",
		TimeWhite:    baseTime,
		TimeBlack:    baseTime,
		BaseTime:     baseTime,
		Increment:    increment,
		WhiteID:      player1,
		BlackID:      player2,
		Board:        board,
		LastMoveTime: time.Now(),
	}
}

func DatabaseGameToGame(game *database.Game) *Game {
	fen, _ := chess.FEN(game.Fen)
	return &Game{
		ID:           game.ID,
		Result:       game.Result,
		BaseTime:     time.Duration(game.BaseTime) * time.Second,
		Increment:    time.Duration(game.Increment) * time.Second,
		TimeBlack:    time.Duration(game.BaseTime) * time.Second,
		TimeWhite:    time.Duration(game.BaseTime) * time.Second,
		WhiteID:      *game.WhiteID,
		BlackID:      *game.BlackID,
		Board:        chess.NewGame(fen),
		GameLength:   game.GameLength,
		LastMoveTime: time.Now(),
	}
}

func (g *Game) MakeMove(move string) (string, string) {
	moveTime := time.Since(g.LastMoveTime)

	if g.Board.Position().Turn() == chess.White {
		if moveTime > g.TimeWhite {
			g.TimeWhite = 0
			g.Result = "0-1"
			return "White Time Out", "0-1"
		} else {
			g.TimeWhite -= moveTime
			g.TimeWhite += g.Increment
		}
	}
	if g.Board.Position().Turn() == chess.Black {
		if moveTime > g.TimeBlack {
			g.TimeBlack = 0
			g.Result = "1-0"
			return "Black Time Out", "1-0"
		} else {
			g.TimeBlack -= moveTime
			g.TimeBlack += g.Increment
		}
	}

	g.LastMoveTime = time.Now()

	if err := g.Board.MoveStr(move); err != nil {
		log.Println(err)
		return "error making move", ""
	}

	g.DrawOfferedBy = 0
	err := g.Board.Draw(chess.ThreefoldRepetition)
	if err != nil {
		g.Board.Draw(chess.FiftyMoveRule)
	}

	if g.Board.Outcome() != "*" {
		g.Result = string(g.Board.Outcome())
		return g.Board.Method().String(), string(g.Board.Outcome())
	}

	return "move successful", ""
}
