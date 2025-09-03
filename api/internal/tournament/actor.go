package tournament

import (
	"sync"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
)

type Tournament struct {
	sync.RWMutex
	ID             string
	Name           string
	Players        map[int32]*Player
	StartTime      time.Time
	Duration       int32
	TimeControl    game.TimeControl
	CreatedBy      int32
	Creator        string
	WaitingPlayers []*Player
	Done           chan struct{}
	BerserkAllowed bool
}

func New(id, name string, duration int32, creator string, createdBy, baseTime, increment int32, numPlayers int, berserkAllowed bool) *Tournament {
	timeControl := game.TimeControl{
		BaseTime:  baseTime,
		Increment: increment,
	}
	t := &Tournament{
		ID:             id,
		Players:        make(map[int32]*Player, numPlayers),
		WaitingPlayers: make([]*Player, 0, numPlayers),
		StartTime:      time.Now(),
		Duration:       duration,
		TimeControl:    timeControl,
		CreatedBy:      createdBy,
		Creator:        creator,
		Name:           name,
		Done:           make(chan struct{}),
		BerserkAllowed: berserkAllowed,
	}
	return t
}
