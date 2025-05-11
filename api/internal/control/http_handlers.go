package control

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func (c *Controller) WriteGameInfo(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "gameID")

	// fmt.Println(gameID)
	foundGame, err := c.Queries.GetGameInfo(r.Context(), gameID)
	if err != nil {
		if challenge, exists := c.GameManager.PendingChallenges[gameID]; exists {
			utils.RespondWithJSON(w, http.StatusOK, challenge)
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}
	serverGame, exists := c.GameManager.Games[gameID]
	if !exists {
		//database game response
		moves, err := c.Queries.GetGameMoves(r.Context(), gameID)
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusBadRequest, "error getting game moves")
			return
		}
		var timeBlack, timeWhite int32
		if foundGame.EndTimeLeftWhite != nil {
			timeWhite = *foundGame.EndTimeLeftWhite
		}
		if foundGame.EndTimeLeftBlack != nil {
			timeBlack = *foundGame.EndTimeLeftBlack
		}
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"moves": moves, "game": foundGame, "timeWhite": timeWhite, "timeBlack": timeBlack})

	} else {
		//server game response
		timePassed := time.Since(serverGame.LastMoveTime)
		if serverGame.Board.Position().Turn() == chess.White {
			serverGame.TimeWhite = max(serverGame.TimeWhite-timePassed, 0)
		} else {
			serverGame.TimeBlack = max(serverGame.TimeBlack-timePassed, 0)
		}
		serverGame.LastMoveTime = time.Now()
		timeWhite := int32(serverGame.TimeWhite.Milliseconds())
		timeBlack := int32(serverGame.TimeBlack.Milliseconds())
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"moves": serverGame.Moves, "game": foundGame, "timeWhite": timeWhite, "timeBlack": timeBlack})
	}
}

func (c *Controller) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.MiddlewareSentSession).(database.GetSessionRow)

	c.SocketManager.RemoveClient(session.UserID)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	err := c.Queries.DeleteSession(r.Context(), session.ID)
	if err != nil {
		log.Printf("error deleting session: %v", err)
		return
	}
}
