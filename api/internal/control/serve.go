package control

import (
	"github.com/gorilla/websocket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
	"log"
	"net/http"
)

var webSocketUpgrader = websocket.Upgrader{
	// CheckOrigin:     checkOrigin,
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Controller) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection request")

	session := r.Context().Value(auth.MiddlewareSentSession).(database.Session)
	username, err := c.Queries.GetUsernameByUserID(r.Context(), session.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "username not set or server error")
		return
	}
	// upgrading http to websocket connection

	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create New Client
	client := socket.NewClient(conn, c.SocketManager, session.UserID, *username)
	// Add the newly created client to the manager
	c.SocketManager.AddClient(session.UserID, client)

	// for _, client := range m.clients {
	// 	fmt.Println(client.UserID)
	// }

	// Start the read / write processes
	go client.ReadMessages()
	go client.WriteMessages()
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
