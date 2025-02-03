package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/control"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	// loading env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}

	setupAPI()
}

func setupAPI() {

	// database setup
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("PG_BREW_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbURL)

	if err != nil {
		log.Fatal("can't connect to database")
	}

	// making sure connection closes when the server stops
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal("can't close connection")
		}
	}(conn, ctx)

	queries := database.New(conn)

	authHandler := auth.NewHandler(queries)

	router := routes.NewChiRouter(authHandler)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	//ws setup
	controller := control.NewController(authHandler)

	router.HandleFunc("/ws", controller.SocketManager.ServeWS)

	log.Printf("Server starting on port: %v", portString)

	// server listening
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
