package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
)

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}

func (h *Handler) ValidateSession(ctx context.Context, token string) (database.GetSessionRow, error) {
	session, err := h.queries.GetSession(ctx, token)
	if err != nil {
		return session, errors.New("no such session")
	}
	if time.Now().After(session.ExpiresAt) {
		err := h.queries.DeleteSession(ctx, token)
		if err != nil {
			log.Println("error deleting session", err)
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

func decodeIDToken(idToken string) (*GoogleIDToken, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode ID Token")
	}

	var user GoogleIDToken
	if err := json.Unmarshal(payload, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user info")
	}

	return &user, nil
}
