package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/control"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

// SetUP chi routes

func NewChiRouter(authHandler *auth.Handler, controller *control.Controller) *chi.Mux {
	router := chi.NewRouter()
	router.Use(utils.RecoveryMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// authenticated routes
	router.Group(func(router chi.Router) {
		router.Use(authHandler.AuthMiddleware)
		router.Post("/logout", authHandler.HandleLogout)
		router.Get("/me", authHandler.HandleMe)

		router.HandleFunc("/ws", controller.SocketManager.ServeWS)
	})

	// routes that don't need authentication
	router.Post("/register", authHandler.HandleRegister)
	router.Post("/login", authHandler.HandleLogin)
	router.Post("/meow", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		log.Println(username)
	})

	// game subroute

	gameRouter := chi.NewRouter()
	gameRouter.Use(authHandler.AuthMiddleware)
	gameRouter.Get("/{gameID}", controller.WriteGameInfo)
	gameRouter.Post("/init", controller.InitGame)

	// mounting subrouters
	router.Mount("/game", gameRouter)

	return router
}
