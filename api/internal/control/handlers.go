package control

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
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

// tournament join leave stuff

func handleJoinLeave(c *Controller, event socket.Event, client *socket.Client) error {
	payload, err := json.Marshal(map[string]any{})
	if err != nil {
		return err
	}
	e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
	var tidP TournamentIDPayload
	if err = json.Unmarshal(event.Payload, &tidP); err != nil {
		client.Send(e)
		return err
	}
	if tidP.TournamentID == "" {
		client.Send(e)
		return errors.New("tournament ID can't be empty")
	}
	t, ok := c.tournamentManager.getTournament(tidP.TournamentID)
	if !ok {
		err = handleJoinLeaveBeforeTournament(c, e, client, tidP.TournamentID)
		return err
	} else {
		err = handleJoinLeaveDuringTournament(c, e, client, t)
		return err
	}
}

func handleJoinLeaveBeforeTournament(c *Controller, e socket.Event, client *socket.Client, tournamentID string) error {
	status, err := c.queries.GetTournamentStatus(context.Background(), tournamentID)
	if err != nil {
		client.Send(e)
		return fmt.Errorf("error getting tournament status for tournament: %s : %w", tournamentID, err)
	}
	if status == 2 {
		client.Send(e)
		return nil
	}
	_, err = c.queries.GetTournamentPlayer(context.Background(), database.GetTournamentPlayerParams{
		PlayerID:     client.UserID,
		TournamentID: tournamentID,
	})
	if err != nil {
		err = c.queries.InsertTournamentPlayer(context.Background(), database.InsertTournamentPlayerParams{
			PlayerID:     client.UserID,
			TournamentID: tournamentID,
		})
		if err != nil {
			client.Send(e)
			return fmt.Errorf("error inserting player to tournament:%s : %w", tournamentID, err)
		}
		rating, err := c.queries.GetUserRating(context.Background(), client.UserID)
		if err != nil {
			client.Send(e)
			return fmt.Errorf("error getting user: %d rating: %w", client.UserID, err)
		}
		payload, err := json.Marshal(map[string]any{"player": map[string]any{"ID": client.UserID, "Score": 0, "Username": client.Username, "Rating": rating}})
		if err != nil {
			client.Send(e)
			return err
		}
		e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
		c.socketManager.BroadcastToRoom(e, tournamentID)
	} else {
		err := c.queries.DeleteTournamentPlayer(context.Background(), database.DeleteTournamentPlayerParams{
			TournamentID: tournamentID,
			PlayerID:     client.UserID,
		})
		if err != nil {
			client.Send(e)
			return fmt.Errorf("error deleting player from tournament:%s : %w", tournamentID, err)
		}
		payload, err := json.Marshal(map[string]any{"id": client.UserID})
		if err != nil {
			return err
		}
		e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
		c.socketManager.BroadcastToRoom(e, tournamentID)
	}
	return nil
}

func handleJoinLeaveDuringTournament(c *Controller, e socket.Event, client *socket.Client, t *tournament.Tournament) error {
	p, exists := t.Players[client.UserID]
	if exists {
		p.IsActive = !p.IsActive
		payload, err := json.Marshal(map[string]any{"player": map[string]any{"ID": client.UserID, "IsActive": p.IsActive}})
		if err != nil {
			client.Send(e)
			return err
		}
		e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
		c.socketManager.BroadcastToRoom(e, t.ID)
	} else {
		err := c.queries.InsertTournamentPlayer(context.Background(), database.InsertTournamentPlayerParams{
			PlayerID:     client.UserID,
			TournamentID: t.ID,
		})
		if err != nil {
			client.Send(e)
			return fmt.Errorf("error inserting player to tournament:%s : %w", t.ID, err)
		}
		rating, err := c.queries.GetUserRating(context.Background(), client.UserID)
		if err != nil {
			client.Send(e)
			return err
		}
		p := tournament.NewPlayer(client.UserID, rating, true)
		t.Players[client.UserID] = p
		t.WaitingPlayers = append(t.WaitingPlayers, p)

		payload, err := json.Marshal(map[string]any{"player": map[string]any{"ID": client.UserID, "Score": p.Score, "Username": client.Username, "Rating": rating, "IsActive": p.IsActive}})
		if err != nil {
			client.Send(e)
			return err
		}
		e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
		c.socketManager.BroadcastToRoom(e, t.ID)
	}
	return nil
}
