package game

import "sync"

type Manager struct {
	sync.Mutex
	Games            map[int32]*Game
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
		Games: make(map[int32]*Game),
	}
}
