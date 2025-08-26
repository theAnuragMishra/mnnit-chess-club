package game

import (
	"sync"
)

type Manager struct {
	sync.RWMutex
	games             map[string]*Game
	pendingChallenges map[string]Challenge
	rematchCache      map[string]*RematchInfo
}

func NewManager() *Manager {
	return &Manager{
		pendingChallenges: make(map[string]Challenge),
		games:             make(map[string]*Game),
		rematchCache:      make(map[string]*RematchInfo),
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

func (m *Manager) AddChallenge(id string, c Challenge) {
	m.Lock()
	m.pendingChallenges[id] = c
	m.Unlock()
}

func (m *Manager) RemoveChallenge(id string) {
	m.Lock()
	delete(m.pendingChallenges, id)
	m.Unlock()
}

func (m *Manager) GetChallengeByID(id string) (Challenge, bool) {
	m.RLock()
	c, exists := m.pendingChallenges[id]
	m.RUnlock()
	return c, exists
}

func (m *Manager) AddRematch(id string, info *RematchInfo) {
	m.Lock()
	m.rematchCache[id] = info
	m.Unlock()
}
func (m *Manager) RemoveRematch(id string) {
	m.Lock()
	delete(m.rematchCache, id)
	m.Unlock()
}

func (m *Manager) GetRematchByID(id string) (*RematchInfo, bool) {
	m.RLock()
	info, exists := m.rematchCache[id]
	m.RUnlock()
	return info, exists
}
