package socket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/config"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
)

var webSocketUpgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
	// CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clientsMu sync.RWMutex
	clients   ClientList
	RoomsMu   sync.RWMutex
	Rooms     map[string]map[*Client]bool
	OnMessage func(event Event, client *Client) error
}

func NewManager(onMessage func(event Event, client *Client) error) *Manager {
	m := &Manager{
		clients:   make(ClientList),
		OnMessage: onMessage,
		Rooms:     make(map[string]map[*Client]bool),
	}
	return m
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection request")

	session := r.Context().Value(auth.MiddlewareSentSession).(database.GetSessionRow)

	if session.Username == nil {
		return
	}

	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, m, session.UserID, *session.Username)
	m.addClient(client)

	// for _, client := range m.clients {
	// 	fmt.Println(client.UserID)
	// }

	// Start the read / write processes
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()
	log.Println("adding client", client)
	if m.clients[client.UserID] == nil {
		m.clients[client.UserID] = make(map[*Client]struct{})
	}
	m.clients[client.UserID][client] = struct{}{}
}

func (m *Manager) RemoveClient(client *Client) {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()

	if clients, ok := m.clients[client.UserID]; ok {
		if _, ok := clients[client]; ok {
			err := client.connection.Close()
			if err != nil {
				log.Println("error trying to close connection of ", client, err)
				return
			}
			close(client.egress)
			delete(m.Rooms[client.Room], client)
			delete(m.clients[client.UserID], client)
			if len(m.Rooms[client.Room]) == 0 {
				delete(m.Rooms, client.Room)
			}
			if len(m.clients[client.UserID]) == 0 {
				delete(m.clients, client.UserID)
			}
		}
	}
}

func (m *Manager) RemoveUser(id int32) {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()

	if clients, ok := m.clients[id]; ok {
		for client := range clients {
			err := client.connection.Close()
			if err != nil {
				return
			}
			delete(m.Rooms[client.Room], client)
			if len(m.Rooms[client.Room]) == 0 {
				delete(m.Rooms, client.Room)
			}
		}
	}
	delete(m.clients, id)
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	return origin == config.FrontendURL
}
