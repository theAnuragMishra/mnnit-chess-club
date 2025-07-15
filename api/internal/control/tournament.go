package control

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
	"log"
	"net/http"
	"time"
)

func (c *Controller) CreateTournament(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.MiddlewareSentSession).(database.GetSessionRow)

	var tournamentPayload TournamentPayload
	if err := json.NewDecoder(r.Body).Decode(&tournamentPayload); err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if tournamentPayload.Name == "" || len(tournamentPayload.Name) > 90 {
		utils.RespondWithError(w, http.StatusBadRequest, "Tournament name should have a length between 1 and 90 characters")
		return
	}

	if tournamentPayload.StartTime.Before(time.Now()) {
		utils.RespondWithError(w, http.StatusBadRequest, "Tournament start time cannot be in the past")
		return
	}

	if tournamentPayload.Duration <= 0 || tournamentPayload.Duration > 86400 {
		utils.RespondWithError(w, http.StatusBadRequest, "Tournament duration must be between 1 and 24 hours")
		return
	}

	if tournamentPayload.BaseTime <= 0 || tournamentPayload.BaseTime > 10800 || tournamentPayload.Increment < 0 || tournamentPayload.Increment > 180 {
		utils.RespondWithError(w, http.StatusBadRequest, "time control not allowed")
		return
	}

	id, err := c.generateUniqueTournamentID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error creating tournament")
		return
	}

	err = c.Queries.CreateTournament(r.Context(), database.CreateTournamentParams{
		ID:        id,
		Name:      tournamentPayload.Name,
		StartTime: tournamentPayload.StartTime,
		Duration:  tournamentPayload.Duration,
		BaseTime:  tournamentPayload.BaseTime,
		Increment: tournamentPayload.Increment,
		CreatedBy: &session.UserID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Tournament name already in use")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]any{"id": id})
}

func (c *Controller) WriteTournamentInfo(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "tournamentID")
	if tournamentID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid tournament ID")
		return
	}
	dbPlayers, err2 := c.Queries.GetTournamentPlayers(r.Context(), tournamentID)
	serverTournament, exists := c.TournamentManager.Tournaments[tournamentID]
	if exists {
		players := make([]any, len(dbPlayers))
		for i, player := range dbPlayers {
			players[i] = map[string]any{
				"Score":    serverTournament.Players[player.ID].Score,
				"ID":       player.ID,
				"IsActive": serverTournament.Players[player.ID].IsActive,
				"Username": player.Username,
				"Rating":   player.Rating,
			}
		}
		//log.Println(players)
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"name":      serverTournament.Name,
			"players":   players,
			"startTime": serverTournament.StartTime,
			"duration":  serverTournament.Duration,
			"baseTime":  serverTournament.TimeControl.BaseTime,
			"increment": serverTournament.TimeControl.Increment,
			"createdBy": serverTournament.CreatedBy,
			"creator":   serverTournament.Creator,
			"ongoing":   true,
		})
		return
	}
	tournamentInfo, err1 := c.Queries.GetTournamentInfo(r.Context(), tournamentID)

	if err1 != nil || err2 != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid tournament ID")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]any{
		"name":      tournamentInfo.Name,
		"players":   dbPlayers,
		"startTime": tournamentInfo.StartTime,
		"duration":  tournamentInfo.Duration,
		"baseTime":  tournamentInfo.BaseTime,
		"increment": tournamentInfo.Increment,
		"createdBy": tournamentInfo.CreatedBy,
		"creator":   tournamentInfo.Username,
		"ongoing":   false,
	})
}

func (c *Controller) StartTournament(w http.ResponseWriter, r *http.Request) {
	var tournamentID TournamentIDPayload
	if json.NewDecoder(r.Body).Decode(&tournamentID) != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if tournamentID.TournamentID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid tournament ID")
		return
	}
	tournamentInfo, err := c.Queries.GetTournamentInfo(r.Context(), tournamentID.TournamentID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error getting tournament info")
		return
	}

	err = c.Queries.UpdateTournamentStartTime(context.Background(), tournamentInfo.ID)
	if err != nil {
		log.Println("Error updating tournament start time", err)
	}

	players, err := c.Queries.GetTournamentPlayers(r.Context(), tournamentID.TournamentID)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error getting tournament players")
		return
	}

	c.TournamentManager.Lock()
	createdTournament := tournament.NewTournament(tournamentInfo.ID, tournamentInfo.Name, tournamentInfo.Duration, *tournamentInfo.Username, *tournamentInfo.CreatedBy, tournamentInfo.BaseTime, tournamentInfo.Increment, len(players))
	c.TournamentManager.Tournaments[tournamentInfo.ID] = createdTournament
	c.TournamentManager.Unlock()

	createdTournament.Lock()
	for i, player := range players {
		p := tournament.NewPlayer(player.ID, player.Rating)
		createdTournament.Players[player.ID] = p
		createdTournament.WaitingPlayers[i] = p
	}
	createdTournament.Unlock()
	//send refresh event to all the players on the tournament page
	payload := map[string]any{"ID": tournamentInfo.ID, "Type": "tournament"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
	c.SocketManager.BroadcastToRoom(e, tournamentInfo.ID)
	time.Sleep(time.Second * 10)

	time.AfterFunc(time.Duration(tournamentInfo.Duration)*time.Second, func() { c.EndTournament(createdTournament) })

	c.RunPairingCycle(createdTournament, true)
	c.StartPairingCycle(createdTournament, time.Second*20)

	utils.RespondWithJSON(w, http.StatusOK, "")
}

func JoinLeaveTournament(c *Controller, event socket.Event, client *socket.Client) error {
	rawPayload, err := json.Marshal(map[string]any{})
	if err != nil {
		return err
	}
	e := socket.Event{Type: "jl_response", Payload: json.RawMessage(rawPayload)}
	var payload TournamentIDPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		client.Send(e)
		return err
	}
	if payload.TournamentID == "" {
		client.Send(e)
		return errors.New("invalid tournament id")
	}
	t, ok := c.TournamentManager.Tournaments[payload.TournamentID]
	if !ok {
		startTimeAndDuration, err := c.Queries.GetTournamentStartTime(context.Background(), payload.TournamentID)
		if err != nil {
			client.Send(e)
			return err
		}
		if time.Now().After(startTimeAndDuration.StartTime.Add(time.Duration(startTimeAndDuration.Duration) * time.Second)) {
			client.Send(e)
			return errors.New("tournament ended")
		}

		tpId, err := c.Queries.GetTournamentPlayer(context.Background(), database.GetTournamentPlayerParams{
			PlayerID:     client.UserID,
			TournamentID: payload.TournamentID,
		})

		if err != nil {

			err := c.Queries.InsertTournamentPlayer(context.Background(), database.InsertTournamentPlayerParams{
				PlayerID:     client.UserID,
				TournamentID: payload.TournamentID,
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
			rawPayloadd, err := json.Marshal(map[string]any{"player": map[string]any{"ID": client.UserID, "Score": 0,
				"Username": client.Username, "Rating": rating,
			}})

			ee := socket.Event{Type: "jl_response", Payload: json.RawMessage(rawPayloadd)}
			c.SocketManager.BroadcastToRoom(ee, payload.TournamentID)
			return nil
		} else {
			err := c.Queries.DeleteTournamentPlayer(context.Background(), tpId)
			if err != nil {
				client.Send(e)
				return err
			}
			rawPayloadd, err := json.Marshal(map[string]any{"id": client.UserID})
			if err != nil {
				return err
			}
			ee := socket.Event{Type: "jl_response", Payload: json.RawMessage(rawPayloadd)}
			c.SocketManager.BroadcastToRoom(ee, payload.TournamentID)
			return nil
		}
	}
	p, exists := t.Players[client.UserID]
	if exists {
		p.Lock()
		p.IsActive = !p.IsActive
		p.Unlock()
		client.Send(e)
	} else {
		err := c.Queries.InsertTournamentPlayer(context.Background(), database.InsertTournamentPlayerParams{
			PlayerID:     client.UserID,
			TournamentID: payload.TournamentID,
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
		player := tournament.NewPlayer(client.UserID, rating)
		t.Lock()
		t.Players[client.UserID] = player
		t.Unlock()
		rawPayloadd, err := json.Marshal(map[string]any{"player": map[string]any{"ID": client.UserID, "Score": player.Score,
			"IsActive": player.IsActive,
			"Username": client.Username,
			"Rating":   player.Rating}})
		if err != nil {
			return err
		}
		ee := socket.Event{Type: "jl_response", Payload: json.RawMessage(rawPayloadd)}
		c.SocketManager.BroadcastToRoom(ee, payload.TournamentID)
	}
	return nil
}

func (c *Controller) WriteUpcomingTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := c.Queries.GetUpcomingTournaments(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting upcoming tournaments")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, tournaments)
}
