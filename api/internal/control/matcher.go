package control

import (
	"sync"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

var timeControls = []game.TimeControl{
	{BaseTime: 60, Increment: 0},
	{BaseTime: 60, Increment: 1},
	{BaseTime: 120, Increment: 1},
	{BaseTime: 180, Increment: 0},
	{BaseTime: 180, Increment: 2},
	{BaseTime: 300, Increment: 0},
	{BaseTime: 300, Increment: 3},
	{BaseTime: 600, Increment: 0},
	{BaseTime: 600, Increment: 5},
	{BaseTime: 900, Increment: 10},
	{BaseTime: 1800, Increment: 0},
	{BaseTime: 1800, Increment: 20}}

type timeControlQueue struct {
	pendingClients []*socket.Client
	timeControl    game.TimeControl
}

func newTimeControlQueue(timeControl game.TimeControl) *timeControlQueue {
	return &timeControlQueue{
		pendingClients: make([]*socket.Client, 0),
		timeControl:    timeControl,
	}
}

type matcher struct {
	sync.Mutex
	queues      []*timeControlQueue
	clientQueue map[*socket.Client]int
}

func newMatcher() *matcher {
	m := &matcher{
		clientQueue: make(map[*socket.Client]int),
	}
	for i := range 12 {
		m.queues = append(m.queues, newTimeControlQueue(timeControls[i]))
	}
	return m
}

func (m *matcher) handleRequest(client *socket.Client, tc int) (int32, bool) {
	m.Lock()
	defer m.Unlock()
	prev, ok := m.clientQueue[client]
	if ok {
		delete(m.clientQueue, client)
		m.removeClient(client, prev)
		if prev == tc {
			return 0, false
		}
	}
	q := m.queues[tc]
	for i, opp := range q.pendingClients {
		if opp.UserID != client.UserID {
			q.pendingClients = append(q.pendingClients[:i], q.pendingClients[i+1:]...)
			delete(m.clientQueue, opp)
			return opp.UserID, true
		}
	}
	q.pendingClients = append(q.pendingClients, client)
	m.clientQueue[client] = tc
	return 0, false
}

func (m *matcher) removeClient(client *socket.Client, tc int) {
	q := m.queues[tc]
	for i, c := range q.pendingClients {
		if c == client {
			q.pendingClients = append(q.pendingClients[:i], q.pendingClients[i+1:]...)
			return
		}
	}
}
