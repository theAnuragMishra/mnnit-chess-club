package game

import "github.com/google/uuid"

type Manager struct {
	Games           []*Game
	PendingUser     uuid.UUID
	PendingUserName string
	// Users         []string
}

func NewManager() *Manager {
	return &Manager{
		Games: []*Game{},
		// Users:         []string{},
	}
}
