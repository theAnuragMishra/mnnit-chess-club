package socket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[int32]*Client

type Client struct {
	// the websocket connection
	connection *websocket.Conn
	UserID     int32
	UserName   string
	manager    *Manager
	egress     chan Event
	Room       int32
}

func NewClient(conn *websocket.Conn, manager *Manager, userID int32, username string) *Client {
	return &Client{
		UserID:     userID,
		UserName:   username,
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		log.Println("client disconnected ", c.UserID)
		c.manager.RemoveClient(c.UserID)
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

		if err := c.manager.OnMessage(request, c); err != nil {
			log.Printf("error on message %v", err)
		}

	}
}

func (c *Client) WriteMessages() {
	defer func() {
		log.Println("Closing write connection for client:", c.UserID)
		c.manager.RemoveClient(c.UserID)
	}()
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
		}
	}
}

func (c *Client) Send(event Event) {
	c.egress <- event
}

func (m *Manager) BroadcastToRoom(event Event, room int32) {
	m.RLock()
	defer m.RUnlock()
	for client := range m.Rooms[room] {
		select {
		case client.egress <- event:
		default:
			log.Printf("Dropping event for client %s: channel full\n", client.connection.RemoteAddr())
		}
	}
}

func (m *Manager) Broadcast(event Event) {
	m.RLock() // Read lock to safely access the clients map
	defer m.RUnlock()

	for _, client := range m.clients {
		select {
		case client.egress <- event:
		default:
			log.Printf("Dropping event for client %s: channel full\n", client.connection.RemoteAddr())
		}
	}
}
