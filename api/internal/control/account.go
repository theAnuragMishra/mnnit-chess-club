package control

import (
	"log"
	"net/http"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
)

func (c *Controller) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.MiddlewareSentSession).(database.GetSessionRow)

	c.SocketManager.RemoveUser(session.UserID)
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

	err := c.Queries.DeleteSession(r.Context(), session.ID)
	if err != nil {
		log.Printf("error deleting session: %v", err)
		return
	}
}
