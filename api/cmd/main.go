package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/config"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/control"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
)

func main() {
	// loading env
	godotenv.Load()
	setupAPI()
}

func setupAPI() {
	// oauth config setup
	config.Config()
	auth.Config()

	// database setup
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("PG_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("can't connect to database")
	}

	// making sure the connection closes when the server stops
	defer pool.Close()

	queries := database.New(pool)

	authHandler := auth.NewHandler(queries)

	controller := control.NewController(queries)

	router := control.NewChiRouter(authHandler, controller)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	auth.StartStateMapCleaner(10 * time.Minute)

	log.Printf("Server starting on port: %v", portString)

	// server listening
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
