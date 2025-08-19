package game

import (
	"sync"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

type Manager struct {
	sync.RWMutex
	Games             map[string]*Game
	PendingUsers      map[TimeControl]PendingUser
	PendingChallenges map[string]Challenge
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

type PendingUser struct {
	ID     int32
	Client *socket.Client
}

func NewManager() *Manager {
	return &Manager{
		Games:             make(map[string]*Game),
		PendingUsers:      make(map[TimeControl]PendingUser),
		PendingChallenges: make(map[string]Challenge),
	}
}
