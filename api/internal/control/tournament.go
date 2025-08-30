package control

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"
)

type tournamentManager struct {
	sync.RWMutex
	tournaments map[string]*tournament.Tournament
}

func newTournamentManager() *tournamentManager {
	return &tournamentManager{
		tournaments: make(map[string]*tournament.Tournament),
	}
}

func (m *tournamentManager) addTournament(t *tournament.Tournament) {
	m.Lock()
	m.tournaments[t.ID] = t
	m.Unlock()
}

func (m *tournamentManager) removeTournament(id string) {
	m.Lock()
	t, ok := m.tournaments[id]
	if ok {
		t.Done <- struct{}{}
	}
	delete(m.tournaments, id)
	m.Unlock()
}

func (m *tournamentManager) getTournament(id string) (*tournament.Tournament, bool) {
	m.RLock()
	t, exists := m.tournaments[id]
	m.RUnlock()
	return t, exists
}

func (c *Controller) endTournament(id string, players []tournament.EndPlayer) {
	err := c.queries.UpdateTournamentStatus(context.Background(), database.UpdateTournamentStatusParams{
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
		err = c.queries.BatchUpdateTournamentPlayers(context.Background(), database.BatchUpdateTournamentPlayersParams{
			TournamentID: id,
			PlayersInput: inputBytes,
		})
		if err != nil {
			log.Println("error batch updating tournament players", err)
		}
	}

	payload := map[string]any{"ID": id, "Type": "tournament"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
	c.socketManager.BroadcastToRoom(e, id)
	c.tournamentManager.removeTournament(id)
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
		err := c.queries.DeleteTournamentPlayer(context.Background(), client.UserID)
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
		msg := tournament.AddPlayer{
			ID:     client.UserID,
			Rating: rating,
			Reply:  make(chan tournament.Player, 1),
		}
		t.Inbox() <- msg
		player := <-msg.Reply
		payload, err := json.Marshal(map[string]any{"player": map[string]any{"ID": client.UserID, "Score": player.Score, "Username": client.Username, "Rating": rating, "IsActive": player.IsActive}})
		if err != nil {
			client.Send(e)
			return err
		}
		e := socket.Event{Type: "jl_response", Payload: json.RawMessage(payload)}
		c.socketManager.BroadcastToRoom(e, t.ID)
	}
	return nil
}
