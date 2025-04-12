package game

import "sync"

type Manager struct {
	sync.Mutex
	Games        map[int32]*Game
	PendingUsers map[string]int32
	// Users         []string
}

func NewManager() *Manager {
	return &Manager{
		Games:        make(map[int32]*Game),
		PendingUsers: make(map[string]int32),
	}
}
