package game

import "github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"

type Manager struct {
	games         []*Game
	pendingGameId string
	Users         []*socket.Client
}

func NewManager() *Manager {
	return &Manager{
		games:         []*Game{},
		pendingGameId: "",
		Users:         []*socket.Client{},
	}
}
