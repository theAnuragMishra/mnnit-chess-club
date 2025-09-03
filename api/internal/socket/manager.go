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
	sync.RWMutex
	clients            ClientList
	rooms              map[string]map[*Client]bool
	OnMessage          func(event Event, client *Client) error
	OnClientDisconnect func(client *Client)
}

func NewManager(onMessage func(event Event, client *Client) error, onClientDisconnect func(client *Client)) *Manager {
	m := &Manager{
		clients:            make(ClientList),
		OnMessage:          onMessage,
		rooms:              make(map[string]map[*Client]bool),
		OnClientDisconnect: onClientDisconnect,
	}
	return m
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	//log.Println("new connection request")

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
	client.sessionID = session.ID
	m.addClient(client)

	// for _, client := range m.clients {
	// 	fmt.Println(client.UserID)
	// }

	// Start the read / write processes
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	//log.Println("adding client", client)
	if m.clients[client.UserID] == nil {
		m.clients[client.UserID] = make(map[*Client]struct{})
	}
	m.clients[client.UserID][client] = struct{}{}
}

func (m *Manager) RemoveClient(client *Client) {
	if clients, ok := m.clients[client.UserID]; ok {
		if _, ok := clients[client]; ok {
			err := client.connection.Close()
			if err != nil {
				log.Println("error trying to close connection of ", client, err)
				return
			}
			delete(m.rooms[client.room], client)
			delete(m.clients[client.UserID], client)
			if len(m.rooms[client.room]) == 0 {
				delete(m.rooms, client.room)
			}
			if len(m.clients[client.UserID]) == 0 {
				delete(m.clients, client.UserID)
			}
		}
	}
}

func (m *Manager) DisconnectAllClientsOfASession(id int32, sessionID string) {
	m.Lock()
	defer m.Unlock()

	if clients, ok := m.clients[id]; ok {
		for client := range clients {
			if client.sessionID == sessionID {
				m.RemoveClient(client)
			}
		}
	}
}

func (m *Manager) AddClientToRoom(room string, client *Client) {
	m.Lock()
	defer m.Unlock()
	if m.rooms[room] == nil {
		m.rooms[room] = make(map[*Client]bool)
	}
	m.rooms[room][client] = true
}

func (m *Manager) DeleteClientFromRoom(room string, client *Client) {
	m.Lock()
	defer m.Unlock()
	delete(m.rooms[room], client)
	if len(m.rooms[room]) == 0 {
		delete(m.rooms, room)
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	return origin == config.FrontendURL
}
