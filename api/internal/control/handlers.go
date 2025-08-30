package control

import (
	"encoding/json"
	"errors"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func chat(c *Controller, event socket.Event, client *socket.Client) error {
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
	g, exists := c.gameManager.getGameByID(client.Room)
	if !exists {
		// handle game ended message
		c.socketManager.BroadcastToRoom(e, client.Room)
		return nil
	}

	if client.UserID != g.WhiteID && client.UserID != g.BlackID {
		// handle message by non player
		c.socketManager.BroadcastToNonPlayers(e, client.Room, g.WhiteID, g.BlackID)
		return nil
	}
	c.socketManager.SendToUserClientsInARoom(e, client.Room, g.WhiteID)
	c.socketManager.SendToUserClientsInARoom(e, client.Room, g.BlackID)
	return nil
}

func roomChange(c *Controller, event socket.Event, client *socket.Client) error {
	var payload RoomPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}
	if client.Room == payload.RoomID {
		return nil
	}
	c.socketManager.DeleteClientFromRoom(client.Room, client)
	client.Room = payload.RoomID
	c.socketManager.AddClientToRoom(payload.RoomID, client)
	return nil
}

func leaveRoom(c *Controller, _ socket.Event, client *socket.Client) error {
	c.socketManager.DeleteClientFromRoom(client.Room, client)
	client.Room = ""
	return nil
}
