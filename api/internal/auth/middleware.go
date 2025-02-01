package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionTokenCookie, err := r.Cookie("session_token")
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		session, err := h.queries.GetSession(r.Context(), sessionTokenCookie.Value)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Session not found")
			return
		}

		if time.Now().After(session.ExpiresAt) {
			utils.RespondWithError(w, http.StatusUnauthorized, "Session expired")
			err := h.queries.DeleteSession(r.Context(), sessionTokenCookie.Value)
			if err != nil {
				log.Println(err)
			}
			return
		}

		if time.Now().Add(time.Hour * 24 * 15).After(session.ExpiresAt) {
			err := h.queries.UpdateSessionExpiry(r.Context(), database.UpdateSessionExpiryParams{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
				ID:        sessionTokenCookie.Value,
			})
			if err != nil {
				log.Println("error updating session expiry, ", err)
			}
		}

		csrf := r.Header.Get("X-CSRF-Token")

		csrfToken, err := h.queries.GetCSRFTokenBySession(r.Context(), sessionTokenCookie.Value)

		if err != nil || csrf != csrfToken.Token || csrfToken.ExpiresAt.Before(time.Now()) {
			utils.RespondWithError(w, http.StatusUnauthorized, "invalid csrf token")
		}

		newCSRFToken := generateToken(32)

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    newCSRFToken,
			Expires:  time.Now().Add(time.Hour),
			HttpOnly: false,
		})

		err = h.queries.UpdateCSRFToken(r.Context(), database.UpdateCSRFTokenParams{
			Token:     newCSRFToken,
			SessionID: sessionTokenCookie.Value,
		})
		if err != nil {
			log.Println("error updating csrf token, ", err)
		}

		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
