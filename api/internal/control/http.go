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

func (c *Controller) handleLogout(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) writeLeaderBoard(w http.ResponseWriter, r *http.Request) {
	lb, err := c.queries.GetTopN(context.Background(), 10)
	if err != nil {
		log.Println("error getting top n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "error fetching leaderboard")
	}
	utils.RespondWithJSON(w, http.StatusOK, lb)
}

func (c *Controller) updateUsername(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.MiddlewareSentSession).(database.GetSessionRow)

	user, err := c.queries.GetUserByUserID(r.Context(), session.UserID)
	if err != nil {
		//log.Println("error getting user by id", err)
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
		log.Println("error updating username for", session.UserID, err)
		utils.RespondWithError(w, http.StatusBadRequest, "Username already in use")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Username updated")
}

func (c *Controller) writeGames(w http.ResponseWriter, r *http.Request) {
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
		log.Println("error getting player games for ", username, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, games)
}

func (c *Controller) writeProfileInfo(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid username")
		return
	}
	profile, err := c.queries.GetUserPublicInfo(r.Context(), &username)
	if err != nil {
		//log.Println("error getting public user info for", username, err)
		utils.RespondWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	counts, err := c.queries.GetGameNumbers(r.Context(), &username)

	if err != nil {
		log.Println("error getting game numbers for", username, err)
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]any{"CreatedAt": profile.CreatedAt, "AvatarUrl": profile.AvatarUrl, "Rating": profile.Rating, "Rd": profile.Rd, "GameCount": counts.GameCount, "DrawCount": counts.DrawCount, "WinCount": counts.WinCount, "LossCount": counts.LossCount})
}

func (c *Controller) writeGameInfo(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "gameID")
	challenge, exists := c.challengeManager.getChallengeByID(gameID)
	if exists {
		utils.RespondWithJSON(w, http.StatusOK, challenge)
		return
	}
	foundGame, err := c.queries.GetGameInfo(r.Context(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}
	g, exists := c.gameManager.getGameByID(gameID)
	if !exists {
		_, canRematch := c.rematchManager.getRematchByID(gameID)
		moves, err := c.queries.GetGameMoves(r.Context(), gameID)
		if err != nil {
			log.Println("error getting game moves for game", gameID, err)
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
		t, ok := c.tournamentManager.getTournament(g.TournamentID)
		berserkAllowed := false
		if ok && t.BerserkAllowed {
			berserkAllowed = true
		}
		msg := game.GetState{Reply: make(chan game.SnapShot, 1)}
		g.Inbox() <- msg
		snapshot := <-msg.Reply
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"moves": snapshot.Moves, "game": foundGame, "timeWhite": snapshot.TimeWhite, "timeBlack": snapshot.TimeBlack, "canRematch": true, "berserkAllowed": berserkAllowed})
	}
}

func (c *Controller) writeTournamentInfo(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "tournamentID")
	if tournamentID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid tournament ID")
		return
	}
	dbPlayers, err := c.queries.GetTournamentPlayers(r.Context(), tournamentID)
	if err != nil {
		log.Println("error getting players for tournament", tournamentID, err)
		utils.RespondWithError(w, 500, "error getting players")
		return
	}
	st, exists := c.tournamentManager.getTournament(tournamentID)
	if exists {
		snapshot := st.PlayersSnapShot.Load().(map[int32]tournament.Player)
		players := make([]any, len(dbPlayers))
		for i, player := range dbPlayers {
			players[i] = map[string]any{
				"Score":    snapshot[player.ID].Score,
				"ID":       player.ID,
				"IsActive": snapshot[player.ID].IsActive,
				"Username": player.Username,
				"Rating":   player.Rating,
				"Streak":   snapshot[player.ID].Streak,
				"Scores":   snapshot[player.ID].Scores,
			}
		}
		//log.Println(players)
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"name":           st.Name,
			"players":        players,
			"startTime":      st.StartTime,
			"duration":       st.Duration,
			"baseTime":       st.TimeControl.BaseTime,
			"increment":      st.TimeControl.Increment,
			"createdBy":      st.CreatedBy,
			"creator":        st.Creator,
			"status":         1,
			"berserkAllowed": st.BerserkAllowed,
		})
		return
	}
	tournamentInfo, err := c.queries.GetTournamentInfo(r.Context(), tournamentID)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid tournament ID")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]any{
		"name":           tournamentInfo.Name,
		"players":        dbPlayers,
		"startTime":      tournamentInfo.StartTime,
		"duration":       tournamentInfo.Duration,
		"baseTime":       tournamentInfo.BaseTime,
		"increment":      tournamentInfo.Increment,
		"createdBy":      tournamentInfo.CreatedBy,
		"creator":        tournamentInfo.Username,
		"status":         tournamentInfo.Status,
		"berserkAllowed": tournamentInfo.BerserkAllowed,
	})
}

func (c *Controller) writeScheduledTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := c.queries.GetScheduledTournaments(r.Context())
	if err != nil {
		log.Println("error getting scheduled tournaments", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting scheduled tournaments")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, tournaments)
}

func (c *Controller) writeLiveTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := c.queries.GetLiveTournaments(r.Context())
	if err != nil {
		log.Println("error getting live tournaments", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting live tournaments")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, tournaments)
}

func (c *Controller) createTournament(w http.ResponseWriter, r *http.Request) {
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
		ID:             id,
		Name:           tournamentPayload.Name,
		StartTime:      tournamentPayload.StartTime,
		Duration:       tournamentPayload.Duration,
		BaseTime:       tournamentPayload.BaseTime,
		Increment:      tournamentPayload.Increment,
		CreatedBy:      &session.UserID,
		BerserkAllowed: tournamentPayload.BerserkAllowed,
	})

	if err != nil {
		log.Println("error creating tournament", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Error creating tournament")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]any{"id": id})
}

func (c *Controller) startTournament(w http.ResponseWriter, r *http.Request) {
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
		utils.RespondWithError(w, http.StatusBadRequest, "Error getting tournament info")
		return
	}

	err = c.queries.UpdateTournamentStartTime(context.Background(), tournamentInfo.ID)
	if err != nil {
		log.Println("error updating tournament start time", err)
	}

	err = c.queries.UpdateTournamentStatus(context.Background(), database.UpdateTournamentStatusParams{
		Status: 1,
		ID:     tournamentInfo.ID,
	})

	if err != nil {
		log.Println("error updating tournament status", err)
	}

	players, err := c.queries.GetTournamentPlayers(r.Context(), tournamentID.TournamentID)

	if err != nil {
		log.Println("error getting tournament players", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Error getting tournament players")
		return
	}

	go func() {
		initialPlayers := make(map[int32]*tournament.Player)
		for _, player := range players {
			p := tournament.NewPlayer(player.ID, player.Rating, c.socketManager.IsUserInARoom(tournamentInfo.ID, player.ID))
			initialPlayers[player.ID] = p
		}

		t := tournament.New(tournamentInfo.ID, tournamentInfo.Name, tournamentInfo.Duration, *tournamentInfo.Username, *tournamentInfo.CreatedBy, tournamentInfo.BaseTime, tournamentInfo.Increment, initialPlayers, c.tournamentRecv, tournamentInfo.BerserkAllowed)
		c.tournamentManager.addTournament(t)

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
