package control

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"
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
	g, exists := c.GameManager.GetGameByID(client.Room)
	if !exists {
		// handle game ended message
		c.SocketManager.BroadcastToRoom(e, client.Room)
		return nil
	}

	if client.UserID != g.WhiteID && client.UserID != g.BlackID {
		// handle message by non player
		c.SocketManager.BroadcastToNonPlayers(e, client.Room, g.WhiteID, g.BlackID)
		return nil
	}
	c.SocketManager.SendToUserClientsInARoom(e, client.Room, g.WhiteID)
	c.SocketManager.SendToUserClientsInARoom(e, client.Room, g.BlackID)
	return nil
}

func createChallenge(c *Controller, event socket.Event, client *socket.Client) error {
	var timeControl game.TimeControl
	if err := json.Unmarshal(event.Payload, &timeControl); err != nil {
		return err
	}
	//log.Println(timeControl)
	if timeControl.BaseTime <= 0 || timeControl.BaseTime > 10800 || timeControl.Increment < 0 || timeControl.Increment > 180 {
		return errors.New("invalid time control")
	}
	id, err := c.generateUniqueGameID()
	if err != nil {
		return err
	}
	c.GameManager.AddChallenge(id, game.Challenge{
		TimeControl:     timeControl,
		Creator:         client.UserID,
		CreatorUsername: client.Username,
	})
	payload := map[string]any{"ID": id, "Type": "game"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}
	client.Send(e)
	return nil
}

func acceptChallenge(c *Controller, event socket.Event, client *socket.Client) error {
	var acceptChallengePayload GameIDPayload
	if err := json.Unmarshal(event.Payload, &acceptChallengePayload); err != nil {
		return err
	}
	challenge, exists := c.GameManager.GetChallengeByID(acceptChallengePayload.GameID)
	if !exists {
		return nil
	}
	if client.UserID == challenge.Creator {
		return nil
	}
	rating1, err := c.Queries.GetUserRating(context.Background(), challenge.Creator)
	if err != nil {
		return err
	}
	rating2, err := c.Queries.GetUserRating(context.Background(), client.UserID)
	if err != nil {
		return err
	}
	g := game.New(acceptChallengePayload.GameID, time.Duration(challenge.TimeControl.BaseTime)*time.Second, time.Duration(challenge.TimeControl.Increment)*time.Second, challenge.Creator, client.UserID, "", c.gameRecv)
	c.GameManager.AddGame(g)
	err = c.Queries.CreateGame(context.Background(), database.CreateGameParams{
		ID:           g.ID,
		BaseTime:     challenge.TimeControl.BaseTime,
		Increment:    challenge.TimeControl.Increment,
		WhiteID:      &challenge.Creator,
		BlackID:      &client.UserID,
		RatingW:      int32(rating1),
		RatingB:      int32(rating2),
		TournamentID: nil,
	})
	if err != nil {
		return err
	}
	payload := map[string]any{"ID": acceptChallengePayload.GameID, "Type": "game"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
	c.SocketManager.BroadcastToRoom(e, acceptChallengePayload.GameID)
	c.GameManager.RemoveChallenge(acceptChallengePayload.GameID)
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
	c.SocketManager.DeleteClientFromRoom(client.Room, client)
	client.Room = payload.RoomID
	c.SocketManager.AddClientToRoom(payload.RoomID, client)

	t, exists := c.TournamentManager.GetTournament(client.Room)
	if exists {
		if c.SocketManager.IsUserInARoom(client.Room, client.UserID) {
			t.Inbox() <- tournament.UpdatePlayerConnectionStatus{
				ID:        client.UserID,
				Connected: true,
			}
		}
	}
	return nil
}

func leaveRoom(c *Controller, _ socket.Event, client *socket.Client) error {
	room := client.Room
	client.Room = ""
	c.SocketManager.DeleteClientFromRoom(room, client)

	t, exists := c.TournamentManager.GetTournament(room)
	if exists {
		if !c.SocketManager.IsUserInARoom(room, client.UserID) {
			t.Inbox() <- tournament.UpdatePlayerConnectionStatus{
				ID:        client.UserID,
				Connected: false,
			}
		}
	}
	return nil
}
