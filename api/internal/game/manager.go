package game

import "sync"

type Manager struct {
	sync.Mutex
	Games             map[string]*Game
	PendingUsers      map[TimeControl]int32
	PendingChallenges map[string]Challenge
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

func NewManager() *Manager {
	return &Manager{
		Games:             make(map[string]*Game),
		PendingUsers:      make(map[TimeControl]int32),
		PendingChallenges: make(map[string]Challenge),
	}
}
