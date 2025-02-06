package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

var webSocketUpgrader = websocket.Upgrader{
	// CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients     ClientList
	authHandler *auth.Handler
	sync.RWMutex

	OnMessage func(event Event, client *Client) error

	// handlers map[string]EventHandler
}

func NewManager(onMessage func(event Event, client *Client) error, authHandler *auth.Handler) *Manager {
	m := &Manager{
		clients:     make(ClientList),
		OnMessage:   onMessage,
		authHandler: authHandler,
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

	// authentication
	sessionTokenCookie, err := r.Cookie("session_token")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	session, err := m.authHandler.ValidateSession(r.Context(), sessionTokenCookie.Value)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Session expired")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionTokenCookie.Value,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	fmt.Println(session)

	// upgrading http to websocket connection

	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create New Client
	client := NewClient(conn, m, session.UserID)
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
	fmt.Println("adding client")
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
