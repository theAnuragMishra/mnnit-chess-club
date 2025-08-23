package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"
)

func (c *Controller) endTournament(id string, players []tournament.EndPlayer) {
	err := c.Queries.UpdateTournamentStatus(context.Background(), database.UpdateTournamentStatusParams{
		Status: 2,
		ID:     id,
	})
	if err != nil {
		log.Println("error updating tournament status", err)
	}

	inputBytes, err := json.Marshal(players)
	if err != nil {
		log.Println(err)
	} else {
		err = c.Queries.BatchUpdateTournamentPlayers(context.Background(), database.BatchUpdateTournamentPlayersParams{
			TournamentID: id,
			PlayersInput: inputBytes,
		})
		if err != nil {
			log.Println(err)
		}
	}

	payload := map[string]any{"ID": id, "Type": "tournament"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
	c.SocketManager.BroadcastToRoom(e, id)
	c.TournamentManager.RemoveTournament(id)
}

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
	t, ok := c.TournamentManager.GetTournament(tidP.TournamentID)
	if !ok {
		err = handleJoinLeaveBeforeTournament(c, e, client, tidP.TournamentID)
		return err
	} else {
		err = handleJoinLeaveDuringTournament(c, e, client, t)
		return err
	}
}

func handleJoinLeaveBeforeTournament(c *Controller, e socket.Event, client *socket.Client, tournamentID string) error {
	status, err := c.Queries.GetTournamentStatus(context.Background(), tournamentID)
	if err != nil {
		client.Send(e)
		return err
	}
	if status == 2 {
		client.Send(e)
		return errors.New("tournament ended")
	}
	_, err = c.Queries.GetTournamentPlayer(context.Background(), database.GetTournamentPlayerParams{
		PlayerID:     client.UserID,
		TournamentID: tournamentID,
	})
	if err != nil {
		err = c.Queries.InsertTournamentPlayer(context.Background(), database.InsertTournamentPlayerParams{
			PlayerID:     client.UserID,
			TournamentID: tournamentID,
		})
		if err != nil {
			client.Send(e)
			return err
		}
		rating, err := c.Queries.GetUserRating(context.Background(), client.UserID)
		if err != nil {
			client.Send(e)
			return err
		}
		payload, err := json.Marshal(map[string]any{"player": map[string]any{"ID": client.UserID, "Score": 0, "Username": client.Username, "Rating": rating}})
		if err != nil {
			client.Send(e)
			return err
		}
		e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
		c.SocketManager.BroadcastToRoom(e, tournamentID)
	} else {
		err := c.Queries.DeleteTournamentPlayer(context.Background(), client.UserID)
		if err != nil {
			client.Send(e)
			return err
		}
		payload, err := json.Marshal(map[string]any{"id": client.UserID})
		if err != nil {
			return err
		}
		e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
		c.SocketManager.BroadcastToRoom(e, tournamentID)
	}
	return nil
}

func handleJoinLeaveDuringTournament(c *Controller, e socket.Event, client *socket.Client, t *tournament.Tournament) error {
	msg := tournament.CheckIfPlayerExists{
		ID:    client.UserID,
		Reply: make(chan bool, 1),
	}
	t.Inbox() <- msg
	exists := <-msg.Reply
	if exists {
		msg := tournament.TogglePlayerActiveMsg{
			ID:    client.UserID,
			Reply: make(chan bool, 1),
		}
		t.Inbox() <- msg
		active := <-msg.Reply
		payload, err := json.Marshal(map[string]any{"player": map[string]any{"ID": client.UserID, "IsActive": active}})
		if err != nil {
			client.Send(e)
			return err
		}
		e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
		c.SocketManager.BroadcastToRoom(e, t.Id)
	} else {
		err := c.Queries.InsertTournamentPlayer(context.Background(), database.InsertTournamentPlayerParams{
			PlayerID:     client.UserID,
			TournamentID: t.Id,
		})
		if err != nil {
			client.Send(e)
			return err
		}
		rating, err := c.Queries.GetUserRating(context.Background(), client.UserID)
		if err != nil {
			client.Send(e)
			return err
		}
		msg := tournament.AddPlayer{
			ID:     client.UserID,
			Rating: rating,
			Reply:  make(chan *tournament.Player, 1),
		}
		t.Inbox() <- msg
		player := <-msg.Reply
		payload, err := json.Marshal(map[string]any{"player": map[string]any{"ID": client.UserID, "Score": player.Score, "Username": client.Username, "Rating": rating, "IsActive": player.IsActive}})
		if err != nil {
			client.Send(e)
			return err
		}
		e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
		c.SocketManager.BroadcastToRoom(e, t.Id)
	}
	return nil
}
