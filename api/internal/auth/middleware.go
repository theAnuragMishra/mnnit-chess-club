package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

type contextKey string

const MiddlewareSentSession contextKey = "session"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionTokenCookie, err := r.Cookie("session_token")
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// fmt.Println("sessionTokenCookie", sessionTokenCookie)

		session, err := h.ValidateSession(r.Context(), sessionTokenCookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: true,
				Path:     "/",
			})

			// also set csrf cookie

			utils.RespondWithError(w, http.StatusUnauthorized, "Session expired")
			return
		}
		if time.Now().Add(time.Hour * 24 * 15).After(session.ExpiresAt) {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    sessionTokenCookie.Value,
				Expires:  time.Now().Add(time.Hour * 24 * 30),
				HttpOnly: true,
				Path:     "/",
			})
		}

		// setting session context for use by the handler
		ctx := context.WithValue(r.Context(), MiddlewareSentSession, session)

		// fmt.Println("passed middleware check")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
