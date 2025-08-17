package control

import (
	"context"
	"log"
	"net/http"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func (c *Controller) WriteLeaderBoard(w http.ResponseWriter, r *http.Request) {
	lb, err := c.Queries.GetTopN(context.Background(), 10)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "error fetching leaderboard")
	}
	utils.RespondWithJSON(w, http.StatusOK, lb)
}
