package socket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
)

var webSocketUpgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
	// CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients ClientList
	sync.RWMutex
	Rooms     map[string]map[*Client]bool
	OnMessage func(event Event, client *Client) error

	// handlers map[string]EventHandler
}

func NewManager(onMessage func(event Event, client *Client) error) *Manager {
	m := &Manager{
		clients:   make(ClientList),
		OnMessage: onMessage,
		Rooms:     make(map[string]map[*Client]bool),
		// handlers: make(map[string]EventHandler),
	}
	// m.setupEventHandlers()
	return m
}

// setting up event handlers
//func (m *Manager) setupEventHandlers() {
//	m.handlers[EventSendMessage] = SendMessage
//}

//func (m *Manager) routeEvent(event Event, c *Client) error {
//	if handler, ok := m.handlers[event.Type]; ok {
//		if err := handler(event, c); err != nil {
//			return err
//		}
//		return nil
//	} else {
//		return errors.New("there is no event of this type")
//	}
//}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection request")

	session := r.Context().Value(auth.MiddlewareSentSession).(database.GetSessionRow)

	if session.Username == nil {
		return
	}

	// upgrading http to websocket connection

	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create New Client
	client := NewClient(conn, m, session.UserID, *session.Username)
	// Add the newly created client to the manager
	m.addClient(client)

	// for _, client := range m.clients {
	// 	fmt.Println(client.UserID)
	// }

	// Start the read / write processes
	go client.readMessages()
	go client.writeMessages()
}

// addClient will add clients to our clientList
func (m *Manager) addClient(client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()
	log.Println("adding client")
	// Add Client
	if m.clients[client.UserID] == nil {
		m.clients[client.UserID] = make(map[*Client]struct{})
	}
	m.clients[client.UserID][client] = struct{}{}
}

// RemoveClient will remove the client and clean up
func (m *Manager) RemoveClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Check if the client exists, then delete it
	if clients, ok := m.clients[client.UserID]; ok {
		if _, ok := clients[client]; ok {
			// close connection
			err := client.connection.Close()
			if err != nil {
				return
			}
			// remove
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

// RemoveUser will remove all the connections of a user
func (m *Manager) RemoveUser(id int32) {
	m.Lock()
	defer m.Unlock()

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
	// Grab the request origin
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:5173":
		return true
	default:
		return false
	}
}
