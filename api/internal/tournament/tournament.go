package tournament

import (
	"sync"
	"time"

	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
)

type Tournament struct {
	sync.RWMutex
	Id             string
	Name           string
	Players        map[int32]*Player
	StartTime      time.Time
	Duration       int32
	TimeControl    game.TimeControl
	CreatedBy      int32
	Creator        string
	WaitingPlayers []*Player
	Done           chan struct{}
}

type Player struct {
	Id              int32
	IsActive        bool
	Score           int32
	Scores          []int16
	Rating          float64
	Streak          int32
	Opponents       map[int32]int16
	LastPlayedColor chess.Color
}

func NewTournament(id, name string, duration int32, creator string, createdBy, baseTime, increment int32, n int) *Tournament {
	timeControl := game.TimeControl{
		BaseTime:  baseTime,
		Increment: increment,
	}
	return &Tournament{
		Id:             id,
		Players:        make(map[int32]*Player),
		WaitingPlayers: make([]*Player, n),
		StartTime:      time.Now(),
		Duration:       duration,
		TimeControl:    timeControl,
		CreatedBy:      createdBy,
		Creator:        creator,
		Name:           name,
		Done:           make(chan struct{}),
	}
}

func NewPlayer(id int32, rating float64) *Player {
	return &Player{
		Id:        id,
		IsActive:  true,
		Score:     0,
		Scores:    make([]int16, 0),
		Rating:    rating,
		Streak:    0,
		Opponents: make(map[int32]int16),
	}
}
