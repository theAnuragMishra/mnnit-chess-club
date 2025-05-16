package control

import (
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"log"
	"net/http"
	"time"
)

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
