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
		log.Println(err)

		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}
	moves, err := c.Queries.GetGameMoves(r.Context(), gameID)
	if err != nil {
		log.Println(err)
	}

	var timeBlack int32
	var timeWhite int32

	serverGame, exists := c.GameManager.Games[gameID]

	if !exists {

		if foundGame.EndTimeLeftWhite != nil {
			timeWhite = *foundGame.EndTimeLeftWhite
		}
		if foundGame.EndTimeLeftBlack != nil {
			timeBlack = *foundGame.EndTimeLeftBlack
		}

	} else {
		timePassed := time.Since(serverGame.LastMoveTime)
		if serverGame.Board.Position().Turn() == chess.White {
			serverGame.TimeWhite = max(serverGame.TimeWhite-timePassed, 0)
		} else {
			serverGame.TimeBlack = max(serverGame.TimeBlack-timePassed, 0)
		}
		serverGame.LastMoveTime = time.Now()

		timeBlack = int32(serverGame.TimeBlack.Milliseconds())
		timeWhite = int32(serverGame.TimeWhite.Milliseconds())
	}

	// fmt.Println(serverGame.Result)

	response := GameResponse{
		Game:      foundGame,
		Moves:     moves,
		TimeBlack: timeBlack,
		TimeWhite: timeWhite,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
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
