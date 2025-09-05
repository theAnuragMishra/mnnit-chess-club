package socket

import (
	"encoding/json"
	"log"
	"sync"

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
		Username:   username,
	}
}

func (c *Client) readMessages() {
	defer func() {
		//log.Println("client disconnected ", c)
		c.manager.OnClientDisconnect(c)
		close(c.egress)
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
	//defer func() {
	//	log.Println("client disconnected", c)
	//}()
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				// if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
				// 	log.Println("connection closed: ", err)
				// }
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("error writing to client: ", c, string(data), err)
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
	c.egress <- event
}

func (m *Manager) BroadcastToRoom(event Event, room string) {
	m.RLock()
	defer m.RUnlock()
	for client := range m.rooms[room] {
		client.egress <- event
	}
}

func (m *Manager) BroadcastToNonPlayers(event Event, room string, player1, player2 int32) {
	m.RLock()
	defer m.RUnlock()
	for client := range m.rooms[room] {
		if client.UserID != player2 && client.UserID != player1 {
			client.egress <- event
		}
	}
}

func (m *Manager) Broadcast(event Event) {
	m.RLock() // Read lock to safely access the clients' map
	defer m.RUnlock()

	for _, clients := range m.clients {
		for client := range clients {
			client.egress <- event
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
				client.egress <- event
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
