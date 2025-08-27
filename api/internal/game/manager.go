package game

import (
	"sync"
)

type Manager struct {
	sync.RWMutex
	games map[string]*Game
}

func NewManager() *Manager {
	return &Manager{
		games: make(map[string]*Game),
	}
}

func (m *Manager) AddGame(g *Game) {
	m.Lock()
	m.games[g.ID] = g
	m.Unlock()
}

func (m *Manager) RemoveGame(id string) {
	m.Lock()
	delete(m.games, id)
	m.Unlock()
}

func (m *Manager) GetGameByID(id string) (*Game, bool) {
	m.RLock()
	g, exists := m.games[id]
	m.RUnlock()
	return g, exists
}
