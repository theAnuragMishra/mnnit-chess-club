package socket

import (
	"log"
	"sync"
)

type Manager struct {
	clients ClientList
	sync.RWMutex

	OnMessage func(event Event, client *Client) error

	// handlers map[string]EventHandler
}

func NewManager(onMessage func(event Event, client *Client) error) *Manager {
	m := &Manager{
		clients:   make(ClientList),
		OnMessage: onMessage,
		// handlers: make(map[string]EventHandler),
	}
	// m.setupEventHandlers()
	return m
}

func (m *Manager) AddClient(id int32, client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()
	log.Println("adding client")
	// Add Client
	m.clients[id] = client
}

func (m *Manager) RemoveClient(id int32) {
	m.Lock()
	defer m.Unlock()

	// Check if Client exists, then delete it
	if client, ok := m.clients[id]; ok {
		// close connection
		err := client.connection.Close()
		if err != nil {
			return
		}
		// remove
		delete(m.clients, id)
	}
}
