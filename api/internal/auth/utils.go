package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}

func (h *Handler) ValidateSession(ctx context.Context, token string) (database.Session, error) {
	session, err := h.queries.GetSession(ctx, token)
	if err != nil {
		return session, errors.New("no such session")
	}
	if time.Now().After(session.ExpiresAt) {
		err := h.queries.DeleteSession(ctx, token)
		if err != nil {
			log.Println(err)
		}
		return session, errors.New("session expired")
	}
	if time.Now().Add(time.Hour * 24 * 15).After(session.ExpiresAt) {
		err := h.queries.UpdateSessionExpiry(ctx, database.UpdateSessionExpiryParams{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
			ID:        token,
		})
		if err != nil {
			log.Println("error updating session expiry, ", err)
		}
	}
	return session, nil
}
