package control

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func (c *Controller) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.MiddlewareSentSession).(database.GetSessionRow)

	c.socketManager.DisconnectAllClientsOfASession(session.UserID, session.ID)
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
		Secure:   auth.CookieCfg.Secure,
		SameSite: auth.CookieCfg.SameSite,
	}
	http.SetCookie(w, cookie)

	err := c.queries.DeleteSession(r.Context(), session.ID)
	if err != nil {
		log.Printf("error deleting session: %v", err)
		return
	}
}

func (c *Controller) WriteLeaderBoard(w http.ResponseWriter, r *http.Request) {
	lb, err := c.queries.GetTopN(context.Background(), 10)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "error fetching leaderboard")
	}
	utils.RespondWithJSON(w, http.StatusOK, lb)
}

func (c *Controller) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.MiddlewareSentSession).(database.GetSessionRow)

	user, err := c.queries.GetUserByUserID(r.Context(), session.UserID)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Internal server error")
		return
	}

	if user.Username != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Username can only be set once")
		return
	}

	var usernamePayload UserNamePayload

	err = json.NewDecoder(r.Body).Decode(&usernamePayload)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = c.queries.UpdateUsername(r.Context(), database.UpdateUsernameParams{
		Username: &usernamePayload.Username,
		ID:       session.UserID,
	})
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Username already in use")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Username updated")
}

func (c *Controller) WriteGames(w http.ResponseWriter, r *http.Request) {
	// log.Println("request received")
	username := chi.URLParam(r, "username")
	page := r.URL.Query().Get("page")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	offSet := int32((pageInt - 1) * 15)
	if offSet < 0 {
		offSet = 0
	}

	games, err := c.queries.GetPlayerGames(r.Context(), database.GetPlayerGamesParams{
		Username: &username,
		Limit:    15,
		Offset:   offSet,
	})
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, games)
}

func (c *Controller) WriteProfileInfo(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid username")
		return
	}
	profile, err := c.queries.GetUserPublicInfo(r.Context(), &username)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	counts, err := c.queries.GetGameNumbers(r.Context(), &username)

	if err != nil {
		log.Println(err)
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]any{"CreatedAt": profile.CreatedAt, "AvatarUrl": profile.AvatarUrl, "Rating": profile.Rating, "Rd": profile.Rd, "GameCount": counts.GameCount, "DrawCount": counts.DrawCount, "WinCount": counts.WinCount, "LossCount": counts.LossCount})
}

func (c *Controller) WriteGameInfo(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "gameID")
	challenge, exists := c.gameManager.GetChallengeByID(gameID)
	if exists {
		utils.RespondWithJSON(w, http.StatusOK, challenge)
		return
	}
	foundGame, err := c.queries.GetGameInfo(r.Context(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}
	g, exists := c.gameManager.GetGameByID(gameID)
	if !exists {
		_, canRematch := c.gameManager.GetRematchByID(gameID)
		moves, err := c.queries.GetGameMoves(r.Context(), gameID)
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError, "error getting game moves")
			return
		}
		var timeBlack, timeWhite int32
		if foundGame.EndTimeLeftWhite != nil {
			timeWhite = *foundGame.EndTimeLeftWhite
		}
		if foundGame.EndTimeLeftBlack != nil {
			timeBlack = *foundGame.EndTimeLeftBlack
		}
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"moves": moves, "game": foundGame, "timeWhite": timeWhite, "timeBlack": timeBlack, "canRematch": canRematch})

	} else {
		// server game response
		msg := game.GetState{Reply: make(chan game.SnapShot, 1)}
		g.Inbox() <- msg
		snapshot := <-msg.Reply
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"moves": snapshot.Moves, "game": foundGame, "timeWhite": snapshot.TimeWhite, "timeBlack": snapshot.TimeBlack, "canRematch": true})
	}
}

func (c *Controller) WriteTournamentInfo(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "tournamentID")
	if tournamentID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid tournament ID")
		return
	}
	dbPlayers, err := c.queries.GetTournamentPlayers(r.Context(), tournamentID)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, 500, "error getting players")
		return
	}
	st, exists := c.tournamentManager.GetTournament(tournamentID)
	if exists {
		msg := tournament.GetState{Reply: make(chan tournament.SnapShot, 1)}
		st.Inbox() <- msg
		snapshot := <-msg.Reply
		players := make([]any, len(dbPlayers))
		for i, player := range dbPlayers {
			players[i] = map[string]any{
				"Score":    snapshot.Players[player.ID].Score,
				"ID":       player.ID,
				"IsActive": snapshot.Players[player.ID].IsActive,
				"Username": player.Username,
				"Rating":   player.Rating,
				"Streak":   snapshot.Players[player.ID].Streak,
				"Scores":   snapshot.Players[player.ID].Scores,
			}
		}
		//log.Println(players)
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"name":      snapshot.Name,
			"players":   players,
			"startTime": snapshot.StartTime,
			"duration":  snapshot.Duration,
			"baseTime":  snapshot.TimeControl.BaseTime,
			"increment": snapshot.TimeControl.Increment,
			"createdBy": snapshot.CreatedBy,
			"creator":   snapshot.Creator,
			"status":    1,
		})
		return
	}
	tournamentInfo, err := c.queries.GetTournamentInfo(r.Context(), tournamentID)

	if err != nil {
		log.Println(err)
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
		"status":    tournamentInfo.Status,
	})
}

func (c *Controller) WriteScheduledTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := c.queries.GetScheduledTournaments(r.Context())
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting scheduled tournaments")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, tournaments)
}

func (c *Controller) WriteLiveTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := c.queries.GetLiveTournaments(r.Context())
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting live tournaments")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, tournaments)
}

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
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Error creating tournament")
		return
	}

	err = c.queries.CreateTournament(r.Context(), database.CreateTournamentParams{
		ID:        id,
		Name:      tournamentPayload.Name,
		StartTime: tournamentPayload.StartTime,
		Duration:  tournamentPayload.Duration,
		BaseTime:  tournamentPayload.BaseTime,
		Increment: tournamentPayload.Increment,
		CreatedBy: &session.UserID,
	})

	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Error creating tournament")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]any{"id": id})
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
	tournamentInfo, err := c.queries.GetTournamentInfo(r.Context(), tournamentID.TournamentID)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Error getting tournament info")
		return
	}

	err = c.queries.UpdateTournamentStartTime(context.Background(), tournamentInfo.ID)
	if err != nil {
		log.Println("Error updating tournament start time", err)
	}

	err = c.queries.UpdateTournamentStatus(context.Background(), database.UpdateTournamentStatusParams{
		Status: 1,
		ID:     tournamentInfo.ID,
	})

	if err != nil {
		log.Println("Error updating tournament status", err)
	}

	players, err := c.queries.GetTournamentPlayers(r.Context(), tournamentID.TournamentID)

	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Error getting tournament players")
		return
	}

	go func() {
		initialPlayers := make(map[int32]*tournament.Player)
		for _, player := range players {
			p := tournament.NewPlayer(player.ID, player.Rating, c.socketManager.IsUserInARoom(tournamentInfo.ID, player.ID))
			initialPlayers[player.ID] = p
		}

		t := tournament.New(tournamentInfo.ID, tournamentInfo.Name, tournamentInfo.Duration, *tournamentInfo.Username, *tournamentInfo.CreatedBy, tournamentInfo.BaseTime, tournamentInfo.Increment, initialPlayers, c.tournamentRecv)
		c.tournamentManager.AddTournament(t)

		//send refresh event to all the players on the tournament page
		payload := map[string]any{"ID": tournamentInfo.ID, "Type": "tournament"}
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			log.Println(err)
		}
		e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
		c.socketManager.BroadcastToRoom(e, tournamentInfo.ID)
	}()
	utils.RespondWithJSON(w, http.StatusOK, "")
}
