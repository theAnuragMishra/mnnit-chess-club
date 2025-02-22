package game

type Manager struct {
	Games            []*Game
	PendingUserName1 string
	PendingUserID1   int32
	PendingUserName2 string
	PendingUserID2   int32
	PendingUserName3 string
	PendingUserID3   int32
	// Users         []string
}

func NewManager() *Manager {
	return &Manager{
		Games: []*Game{},
	}
}
