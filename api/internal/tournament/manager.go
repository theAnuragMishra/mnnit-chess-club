package tournament

import "sync"

type Manager struct {
	sync.RWMutex
	Tournaments map[string]*Tournament
}

func NewManager() *Manager {
	return &Manager{
		Tournaments: make(map[string]*Tournament),
	}
}
