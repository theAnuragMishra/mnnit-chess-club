package game

import (
	"time"

	"github.com/notnil/chess"
)

type Game struct {
	ID            string
	BaseTime      time.Duration
	Increment     time.Duration
	WhiteID       int32
	BlackID       int32
	Board         *chess.Game
	LastMoveTime  time.Time
	TimeWhite     time.Duration
	TimeBlack     time.Duration
	DrawOfferedBy int32
	AbortTimer    *time.Timer
	ClockTimer    *time.Timer
	Moves         []Move
	TournamentID  string
}

type Move struct {
	MoveNotation string
	Orig         string
	Dest         string
	MoveFen      string
	TimeLeft     *int32
}

func NewGame(baseTime time.Duration, increment time.Duration, player1 int32, player2 int32, tournamentID string) *Game {
	board := chess.NewGame()

	return &Game{
		TimeWhite:    baseTime,
		TimeBlack:    baseTime,
		BaseTime:     baseTime,
		Increment:    increment,
		WhiteID:      player1,
		BlackID:      player2,
		Board:        board,
		LastMoveTime: time.Now(),
		TournamentID: tournamentID,
	}
}

func (g *Game) MakeMove(move string) error {
	moveTime := time.Since(g.LastMoveTime)

	if g.Board.Position().Turn() == chess.White {
		g.TimeWhite -= moveTime
		g.TimeWhite += g.Increment
	} else {
		g.TimeBlack -= moveTime
		g.TimeBlack += g.Increment
	}

	g.LastMoveTime = time.Now()

	if err := g.Board.MoveStr(move); err != nil {
		//log.Println(err)
		return err
	}

	g.DrawOfferedBy = 0
	err := g.Board.Draw(chess.ThreefoldRepetition)
	if err != nil {
		g.Board.Draw(chess.FiftyMoveRule)
	}

	return nil
}
