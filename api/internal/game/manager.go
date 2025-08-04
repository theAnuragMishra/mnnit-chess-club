package game

import (
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"sync"
)

type Manager struct {
	sync.Mutex
	Games             map[string]*Game
	PendingUsers      map[TimeControl]PendingUser
	PendingChallenges map[string]Challenge
	Rematches         map[string]Rematch
	// Users         []string
}

type TimeControl struct {
	BaseTime  int32 `json:"baseTime"`
	Increment int32 `json:"increment"`
}

type Challenge struct {
	TimeControl     TimeControl
	Creator         int32
	CreatorUsername string
}

type Rematch struct {
	TimeControl TimeControl
	Creator     int32
	Opponent    int32
}

type PendingUser struct {
	ID     int32
	Client *socket.Client
}

func NewManager() *Manager {
	return &Manager{
		Games:             make(map[string]*Game),
		PendingUsers:      make(map[TimeControl]PendingUser),
		PendingChallenges: make(map[string]Challenge),
		Rematches:         make(map[string]Rematch),
	}
}
