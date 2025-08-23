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

func (m *Manager) AddTournament(t *Tournament) {
	m.Lock()
	m.Tournaments[t.Id] = t
	m.Unlock()
}

func (m *Manager) RemoveTournament(id string) {
	m.Lock()
	delete(m.Tournaments, id)
	m.Unlock()
}

func (m *Manager) GetTournament(id string) (*Tournament, bool) {
	m.RLock()
	t, exists := m.Tournaments[id]
	m.RUnlock()
	return t, exists
}
