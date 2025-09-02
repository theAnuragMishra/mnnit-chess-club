package game

import (
	"sync"
	"time"

	"github.com/notnil/chess"
)

type Game struct {
	sync.RWMutex
	ID            string
	BaseTime      time.Duration
	Increment     time.Duration
	WhiteID       int32
	BlackID       int32
	TournamentID  string
	BerserkWhite  bool
	BerserkBlack  bool
	Result        int
	Board         *chess.Game
	TimeWhite     time.Duration
	TimeBlack     time.Duration
	DrawOfferedBy int32
	Moves         []Move
	RematchOffer  bool
	LastMoveTime  time.Time
	AbortTimer    *time.Timer
	ClockTimer    *time.Timer
	EndCallback   func(g *Game, info EndInfo)
}

func New(id string, baseTime time.Duration, increment time.Duration, player1 int32, player2 int32, tournamentID string, endCallback func(g *Game, info EndInfo)) *Game {
	g := &Game{
		ID:           id,
		BaseTime:     baseTime,
		Increment:    increment,
		Board:        chess.NewGame(),
		TimeBlack:    baseTime,
		TimeWhite:    baseTime,
		LastMoveTime: time.Now(),
		WhiteID:      player1,
		BlackID:      player2,
		TournamentID: tournamentID,
		EndCallback:  endCallback,
	}
	g.setUpAbort()
	return g
}
