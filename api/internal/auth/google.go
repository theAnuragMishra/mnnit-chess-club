package auth

import (
	"context"
	"fmt"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := Config().AuthCodeURL(os.Getenv("OAUTH_STATE"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

type GoogleIDToken struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Google Callback")
	code := r.URL.Query().Get("code")
	token, err := Config().Exchange(context.Background(), code)
	if err != nil {
		log.Println(err, "failed to exchange token")
		http.Redirect(w, r, "http://localhost:5173/error-page", http.StatusTemporaryRedirect)
		return
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Println("No id_token field in oauth2 token")
		http.Redirect(w, r, "http://localhost:5173/error-page", http.StatusTemporaryRedirect)
		return
	}

	googleUser, err := decodeIDToken(idToken)
	if err != nil {
		log.Printf("failed to decode ID Token: %v", err)
		http.Redirect(w, r, "http://localhost:5173/error-page", http.StatusTemporaryRedirect)
		return
	}

	databaseUser, err := h.queries.GetUserByEmail(context.Background(), googleUser.Email)

	if err != nil {
		databaseUser, err = h.queries.CreateUser(r.Context(), database.CreateUserParams{
			UpdatedAt: time.Now(),
			Email:     googleUser.Email,
			AvatarUrl: &googleUser.Picture,
			GoogleID:  googleUser.Sub,
		})

		fmt.Printf("%+v\n", databaseUser)
		fmt.Printf("%+v\n", err)

		if err != nil {
			http.Redirect(w, r, "http://localhost:5173/error-page", http.StatusTemporaryRedirect)
			return

		}
	} else {
		databaseUser, err = h.queries.UpdateUserAvatar(r.Context(), database.UpdateUserAvatarParams{
			AvatarUrl: &googleUser.Picture,
			ID:        databaseUser.ID,
		})

		if err != nil {
			log.Println(err, " failed to update avatar")
			http.Redirect(w, r, "http://localhost:5173/error-page", http.StatusTemporaryRedirect)
			return
		}
	}

	sessionToken := generateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		Path:     "/",
	})
	err = h.queries.CreateSession(r.Context(), database.CreateSessionParams{
		ID:        sessionToken,
		UserID:    databaseUser.ID,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour * 30),
	})
	if err != nil {
		log.Println(err, "failed to create session")
		http.Redirect(w, r, "http://localhost:5173/error-page", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}
