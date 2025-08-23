package game

import (
	"sync"
)

type Manager struct {
	sync.RWMutex
	Games             map[string]*Game
	PendingUsers      map[TimeControl]int32
	PendingChallenges map[string]Challenge
	RematchCache      map[string]*RematchInfo
}

func NewManager() *Manager {
	return &Manager{
		PendingUsers:      make(map[TimeControl]int32),
		PendingChallenges: make(map[string]Challenge),
		Games:             make(map[string]*Game),
		RematchCache:      make(map[string]*RematchInfo),
	}
}

func (m *Manager) AddGame(g *Game) {
	m.Lock()
	m.Games[g.ID] = g
	m.Unlock()
}

func (m *Manager) RemoveGame(id string) {
	m.Lock()
	delete(m.Games, id)
	m.Unlock()
}

func (m *Manager) GetGameByID(id string) (*Game, bool) {
	m.RLock()
	g, exists := m.Games[id]
	m.RUnlock()
	return g, exists
}

func (m *Manager) AddChallenge(id string, c Challenge) {
	m.Lock()
	m.PendingChallenges[id] = c
	m.Unlock()
}

func (m *Manager) RemoveChallenge(id string) {
	m.Lock()
	delete(m.PendingChallenges, id)
	m.Unlock()
}

func (m *Manager) GetChallengeByID(id string) (Challenge, bool) {
	m.RLock()
	c, exists := m.PendingChallenges[id]
	m.RUnlock()
	return c, exists
}

func (m *Manager) AddPendingUser(tc TimeControl, id int32) {
	m.Lock()
	m.PendingUsers[tc] = id
	m.Unlock()
}
func (m *Manager) RemovePendingUser(tc TimeControl) {
	m.Lock()
	delete(m.PendingUsers, tc)
	m.Unlock()
}
func (m *Manager) GetPendingUser(tc TimeControl) (int32, bool) {
	m.RLock()
	user, exists := m.PendingUsers[tc]
	m.RUnlock()
	return user, exists
}

func (m *Manager) AddRematch(id string, info *RematchInfo) {
	m.Lock()
	m.RematchCache[id] = info
	m.Unlock()
}
func (m *Manager) RemoveRematch(id string) {
	m.Lock()
	delete(m.RematchCache, id)
	m.Unlock()
}

func (m *Manager) GetRematchByID(id string) (*RematchInfo, bool) {
	m.RLock()
	info, exists := m.RematchCache[id]
	m.RUnlock()
	return info, exists
}
