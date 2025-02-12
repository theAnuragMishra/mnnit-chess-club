package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
)

var webSocketUpgrader = websocket.Upgrader{
	// CheckOrigin:     checkOrigin,
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

	session := r.Context().Value(auth.MiddlewareSentSession).(database.Session)

	// upgrading http to websocket connection

	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create New Client
	client := NewClient(conn, m, session.UserID)
	// Add the newly created client to the manager
	m.addClient(session.UserID, client)

	// Start the read / write processes
	go client.readMessages()
	go client.writeMessages()
}

// addClient will add clients to our clientList
func (m *Manager) addClient(id uuid.UUID, client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()
	fmt.Println("adding client")
	// Add Client
	m.clients[id] = client
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(id uuid.UUID) {
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
