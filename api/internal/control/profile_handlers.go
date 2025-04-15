package control

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func (c *Controller) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.MiddlewareSentSession).(database.GetSessionRow)

	user, err := c.Queries.GetUserByUserID(r.Context(), session.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Internal server error")
		return
	}

	if user.Username != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Username can only be set once")
	}

	var usernamePayload UserNamePayload

	err = json.NewDecoder(r.Body).Decode(&usernamePayload)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = c.Queries.UpdateUsername(r.Context(), database.UpdateUsernameParams{
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

	games, err := c.Queries.GetPlayerGames(r.Context(), database.GetPlayerGamesParams{
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
	profile, err := c.Queries.GetUserPublicInfo(r.Context(), &username)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	counts, err := c.Queries.GetGameNumbers(r.Context(), &username)

	if err != nil {
		log.Println(err)
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]any{"CreatedAt": profile.CreatedAt, "AvatarUrl": profile.AvatarUrl, "Rating": profile.Rating, "Rd": profile.Rd, "GameCount": counts.GameCount, "DrawCount": counts.DrawCount, "WinCount": counts.WinCount, "LossCount": counts.LossCount})
}
