package auth

import (
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
	"net/http"
)

type Handler struct {
	queries *database.Queries
}

func NewHandler(queries *database.Queries) *Handler {
	return &Handler{queries: queries}
}

// a "/me" route to send userdata on app load

func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(MiddlewareSentSession).(database.GetSessionRow)

	user, err := h.queries.GetUserByUserID(r.Context(), session.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting user")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]any{"username": user.Username, "userID": user.ID, "role": user.Role})
}
