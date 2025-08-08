package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/config"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := generateToken(16)
	stateMap[state] = struct{}{}
	url := oauthCfg.AuthCodeURL(state)
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
	state := r.URL.Query().Get("state")
	_, ok := stateMap[state]
	if !ok {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	token, err := oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		log.Println(err, "failed to exchange token")
		http.Redirect(w, r, config.FrontendURL+"/error-page", http.StatusTemporaryRedirect)
		return
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Println("No id_token field in oauth2 token")
		http.Redirect(w, r, config.FrontendURL+"/error-page", http.StatusTemporaryRedirect)
		return
	}

	googleUser, err := decodeIDToken(idToken)
	if err != nil {
		log.Printf("failed to decode ID Token: %v", err)
		http.Redirect(w, r, config.FrontendURL+"/error-page", http.StatusTemporaryRedirect)
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
		//fmt.Printf("%+v\n", databaseUser)
		//fmt.Printf("%+v\n", err)

		if err != nil {
			http.Redirect(w, r, config.FrontendURL+"/error-page", http.StatusTemporaryRedirect)
			return

		}
	} else if *databaseUser.AvatarUrl != googleUser.Picture {
		err = h.queries.UpdateUserAvatar(r.Context(), database.UpdateUserAvatarParams{
			AvatarUrl: &googleUser.Picture,
			ID:        databaseUser.ID,
			UpdatedAt: time.Now(),
		})

		if err != nil {
			log.Println(err, " failed to update avatar")
		}
	}

	sessionToken := generateToken(32)

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		Path:     "/",
		Secure:   CookieCfg.Secure,
		SameSite: CookieCfg.SameSite,
	}

	http.SetCookie(w, cookie)
	err = h.queries.CreateSession(r.Context(), database.CreateSessionParams{
		ID:        sessionToken,
		UserID:    databaseUser.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30),
	})
	if err != nil {
		log.Println(err, "failed to create session")
		http.Redirect(w, r, config.FrontendURL+"/error-page", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, config.FrontendURL+"/", http.StatusTemporaryRedirect)

}
