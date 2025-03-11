package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

type Handler struct {
	queries *database.Queries
}

func NewHandler(queries *database.Queries) *Handler {
	return &Handler{queries: queries}
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionTokenCookie, err := r.Cookie("session_token")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	err = h.queries.DeleteSession(r.Context(), sessionTokenCookie.Value)
	if err != nil {
		log.Printf("error deleting session: %v", err)
		return
	}
}

// a "/me" route to send userdata on app load

func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(MiddlewareSentSession).(database.Session)

	user, err := h.queries.GetUserByUserID(r.Context(), session.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting user")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"username": user.Username, "userID": user.ID})
}
