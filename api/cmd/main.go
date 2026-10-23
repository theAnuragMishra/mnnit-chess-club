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

	// server port setup
	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8080"
	}

	// database connection setup (supports PG_URL from Docker / local and DATABASE_URL from Render)
	dbURL := os.Getenv("PG_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		log.Fatal("Neither PG_URL nor DATABASE_URL was found in the environment")
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
