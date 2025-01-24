package game

type Manager struct {
	Games         []*Game
	PendingGameId string
	//Users         []string
}

func NewManager() *Manager {
	return &Manager{
		Games: []*Game{},
		//Users:         []string{},
	}
}
