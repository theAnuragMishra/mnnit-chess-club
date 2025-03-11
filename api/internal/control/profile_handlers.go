package control

import (
	"encoding/json"
	"fmt"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
	"net/http"
)

func (c *Controller) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.MiddlewareSentSession).(database.Session)

	var usernamePayload UserNamePayload
	err := json.NewDecoder(r.Body).Decode(&usernamePayload)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = c.Queries.UpdateUsername(r.Context(), database.UpdateUsernameParams{
		Username: &usernamePayload.Username,
		ID:       session.UserID,
	})

	fmt.Println(err)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Username already in use")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Username updated")
}
