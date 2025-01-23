package game

type Manager struct {
	games         []*Game
	pendingGameId string
	Users         []string
}

func NewManager() *Manager {
	return &Manager{
		games:         []*Game{},
		pendingGameId: "",
		Users:         []string{},
	}
}
