package socket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[int32]*Client

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn
	UserID     int32
	// manager is the manager used to manage the client
	manager  *Manager
	egress   chan Event
	Room     string
	Username string
}

// NewClient is used to initialize a new Client with all required values initialized
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
		//if err := c.manager.routeEvent(request, c); err != nil {
		//	log.Println("error handling message: ", err)
		//}

		if err := c.manager.OnMessage(request, c); err != nil {
			log.Printf("error on message %v", err)
		}

	}
}

func (c *Client) writeMessages() {
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

// Send sends an event to the egress channel which is then written to the client by WriteMessage
func (c *Client) Send(event Event) {
	c.egress <- event
}

func (m *Manager) BroadcastToRoom(event Event, room string) {
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

func (m *Manager) BroadcastToNonPlayers(event Event, room string, player1, player2 *Client) {
	m.RLock()
	defer m.RUnlock()
	for client := range m.Rooms[room] {
		if client != player2 && client != player1 {
			select {
			case client.egress <- event:
			default:
				log.Printf("Dropping event for client %s: channel full\n", client.connection.RemoteAddr())
			}
		}
	}
}

// Broadcast sends the event to every client
func (m *Manager) Broadcast(event Event) {
	m.RLock() // Read lock to safely access the clients map
	defer m.RUnlock()

	for _, client := range m.clients {
		select {
		case client.egress <- event:
			// Successfully enqueued the event
		default:
			// Client's channel is full; log and handle as needed
			log.Printf("Dropping event for client %s: channel full\n", client.connection.RemoteAddr())
		}
	}
}
