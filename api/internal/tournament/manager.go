package tournament

import "sync"

type Manager struct {
	sync.RWMutex
	tournaments map[string]*Tournament
}

func NewManager() *Manager {
	return &Manager{
		tournaments: make(map[string]*Tournament),
	}
}

func (m *Manager) AddTournament(t *Tournament) {
	m.Lock()
	m.tournaments[t.ID] = t
	m.Unlock()
}

func (m *Manager) RemoveTournament(id string) {
	m.Lock()
	delete(m.tournaments, id)
	m.Unlock()
}

func (m *Manager) GetTournament(id string) (*Tournament, bool) {
	m.RLock()
	t, exists := m.tournaments[id]
	m.RUnlock()
	return t, exists
}
