package game

type Manager struct {
	Games           []*Game
	PendingUserName string
	PendingUserID   int32
	// Users         []string
}

func NewManager() *Manager {
	return &Manager{
		Games: []*Game{},
		// Users:         []string{},
	}
}
