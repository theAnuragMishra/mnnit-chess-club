package control

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func (c *Controller) WriteProfileInfo(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	profileInfo, err := c.Queries.GetPlayerGames(r.Context(), username)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	}

	utils.RespondWithJSON(w, http.StatusOK, profileInfo)
}

func (c *Controller) WriteGameInfo(w http.ResponseWriter, r *http.Request) {
	gameIDStr, err := strconv.ParseInt(chi.URLParam(r, "gameID"), 10, 32)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
	}
	gameID := int32(gameIDStr)
	// fmt.Println(gameID)
	foundGame, err := c.Queries.GetGameInfo(r.Context(), gameID)
	if err != nil {
		log.Println(err)

		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
	}
	moves, err := c.Queries.GetGameMoves(r.Context(), gameID)
	if err != nil {
		log.Println(err)
	}

	response := GameResponse{
		Game:  foundGame,
		Moves: moves,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
