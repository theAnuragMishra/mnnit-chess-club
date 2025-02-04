package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

type Handler struct {
	queries *database.Queries
}

func NewHandler(queries *database.Queries) *Handler {
	return &Handler{queries: queries}
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var registerRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	username := registerRequest.Username
	password := registerRequest.Password
	// email := r.FormValue("email")
	if len(username) < 4 {
		utils.RespondWithError(w, http.StatusBadRequest, "Username and password must be at least 6 characters long")
		return
	}
	if len(password) < 6 {
		utils.RespondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	err = h.queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	username := loginRequest.Username
	password := loginRequest.Password

	user, err := h.queries.GetUserByUsername(r.Context(), username)

	if err != nil || !checkPasswordHash(password, user.PasswordHash) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: false,
	})

	err = h.queries.CreateSession(r.Context(), database.CreateSessionParams{
		ID:        sessionToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour * 30),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldn't create session")
		return
	}

	err = h.queries.CreateCSRFToken(r.Context(), database.CreateCSRFTokenParams{
		SessionID: sessionToken,
		Token:     csrfToken,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour * 30),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldn't create CSRF token")
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"username": username})
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("into handle logout")

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
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	err = h.queries.DeleteSession(r.Context(), sessionTokenCookie.Value)
	if err != nil {
		log.Printf("error deleting session: %v", err)
		return
	}
}

// a "/me" route to send userdata on app load

func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into handle me")

	session := r.Context().Value(middlewareSentSession).(database.Session)

	fmt.Println("converted session to database.Session type")

	user, err := h.queries.GetUserByUserID(r.Context(), session.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting user")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"username": user.Username})
}
