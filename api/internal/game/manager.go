package game

import "sync"

type Manager struct {
	sync.Mutex
	Games             map[string]*Game
	PendingUsers      map[string]int32
	PendingChallenges map[string]Challenge
	// Users         []string
}

type Challenge struct {
	TimeControl     string
	Creator         int32
	CreatorUsername string
}

func NewManager() *Manager {
	return &Manager{
		Games:             make(map[string]*Game),
		PendingUsers:      make(map[string]int32),
		PendingChallenges: make(map[string]Challenge),
	}
}
