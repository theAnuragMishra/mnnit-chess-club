package tournament

import "sync"

type Manager struct {
	sync.Mutex
	Tournaments map[string]*Tournament
}

func NewManager() *Manager {
	return &Manager{
		Tournaments: make(map[string]*Tournament),
	}
}
