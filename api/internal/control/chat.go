package control

import (
	"encoding/json"
	"errors"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func Chat(c *Controller, event socket.Event, client *socket.Client) error {
	var chat ChatPayload
	if err := json.Unmarshal(event.Payload, &chat); err != nil {
		return errors.New("invalid payload")
	}

	payload, err := json.Marshal(map[string]any{"sender": client.Username, "text": chat.Text})
	if err != nil {
		return errors.New("error marshalling chat payload")
	}

	e := socket.Event{
		Type:    "chat",
		Payload: json.RawMessage(payload),
	}

	foundGame, exists := c.GameManager.Games[client.Room]
	if !exists {
		// handle game ended message
		c.SocketManager.BroadcastToRoom(e, client.Room)
		return nil
	}

	if client.UserID != foundGame.WhiteID && client.UserID != foundGame.BlackID {
		// handle message by non player
		c.SocketManager.BroadcastToNonPlayers(e, client.Room, foundGame.WhiteID, foundGame.BlackID)
		return nil
	}
	c.SocketManager.SendToUserClientsInARoom(e, client.Room, foundGame.WhiteID)
	c.SocketManager.SendToUserClientsInARoom(e, client.Room, foundGame.BlackID)
	return nil
}
