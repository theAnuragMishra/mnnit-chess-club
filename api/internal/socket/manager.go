package socket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var webSocketUpgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients ClientList

	sync.RWMutex

	OnMessage func(event Event, client *Client) error

	//handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients: make(ClientList),
		//handlers: make(map[string]EventHandler),
	}
	//m.setupEventHandlers()
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
	log.Println("new connection")
	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create New Client
	client := NewClient(conn, m)
	// Add the newly created client to the manager
	m.addClient(client)

	// Start the read / write processes
	go client.readMessages()
	go client.writeMessages()
}

// addClient will add clients to our clientList
func (m *Manager) addClient(client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[client] = true
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Check if Client exists, then delete it
	if _, ok := m.clients[client]; ok {
		// close connection
		err := client.connection.Close()
		if err != nil {
			return
		}
		// remove
		delete(m.clients, client)
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
