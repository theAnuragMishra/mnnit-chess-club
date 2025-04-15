package game

import "sync"

type Manager struct {
	sync.Mutex
	Games        map[string]*Game
	PendingUsers map[string]int32
	// Users         []string
}

func NewManager() *Manager {
	return &Manager{
		Games:        make(map[string]*Game),
		PendingUsers: make(map[string]int32),
	}
}
