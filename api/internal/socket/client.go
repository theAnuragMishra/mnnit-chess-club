package socket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[int32]map[*Client]struct{}

// Client is a websocket client, basically a frontend visitor
type Client struct {
	sync.RWMutex
	connection *websocket.Conn
	UserID     int32
	manager    *Manager
	egress     chan Event
	done       chan struct{}
	room       string
	Username   string
	sessionID  string
}

func NewClient(conn *websocket.Conn, manager *Manager, userID int32, username string) *Client {
	return &Client{
		UserID:     userID,
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		done:       make(chan struct{}),
		Username:   username,
	}
}

func (c *Client) readMessages() {
	defer func() {
		close(c.done)
		c.manager.OnClientDisconnect(c)
		c.manager.Lock()
		c.manager.RemoveClient(c)
		c.manager.Unlock()
	}()
	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error", err)
			}
			break
		}
		if c.manager.messageLimit.Take(c.UserID) > 0 {
			c.connection.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Message rate limit exceeded"),
				time.Now().Add(time.Second))
			return
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling event %v", err)
			break
		}
		//if err := c.manager.routeEvent(request, c); err != nil {
		//	log.Println("error handling message: ", err)
		//}

		if err := c.manager.OnMessage(request, c); err != nil {
			log.Printf("error returned by event handler: %v", err)
		}

	}
}

func (c *Client) writeMessages() {
	defer c.connection.Close()
	for {
		select {
		case <-c.done:
			return
		case message := <-c.egress:
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}
			c.connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("error writing to websocket for user %d: %v", c.UserID, err)
				return
			}
		}
	}
}

func (c *Client) Room() string {
	return c.room
}

func (c *Client) ChangeRoom(room string) {
	c.Lock()
	c.room = room
	c.Unlock()
}

func (c *Client) RoomLock() string {
	c.RLock()
	defer c.RUnlock()
	return c.room
}

func (c *Client) Send(event Event) {
	select {
	case c.egress <- event:
	case <-c.done:
	}
}

func (m *Manager) BroadcastToRoom(event Event, room string) {
	m.RLock()
	defer m.RUnlock()
	for client := range m.rooms[room] {
		client.Send(event)
	}
}

func (m *Manager) BroadcastToNonPlayers(event Event, room string, player1, player2 int32) {
	m.RLock()
	defer m.RUnlock()
	for client := range m.rooms[room] {
		if client.UserID != player2 && client.UserID != player1 {
			client.Send(event)
		}
	}
}

func (m *Manager) Broadcast(event Event) {
	m.RLock() // Read lock to safely access the clients' map
	defer m.RUnlock()

	for _, clients := range m.clients {
		for client := range clients {
			client.Send(event)
		}
	}
}

func (m *Manager) SendToUserClientsInARoom(event Event, room string, id int32) {
	m.RLock()
	defer m.RUnlock()
	clients, ok := m.clients[id]
	if ok {
		for client := range clients {
			if client.RoomLock() == room {
				client.Send(event)
			}
		}
	}
}

func (m *Manager) IsUserInARoom(room string, id int32) bool {
	m.RLock()
	defer m.RUnlock()
	clients, ok := m.clients[id]
	if ok {
		for client := range clients {
			if client.RoomLock() == room {
				return true
			}
		}
	}
	return false
}
