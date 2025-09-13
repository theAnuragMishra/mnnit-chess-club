package control

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func NewChiRouter(authHandler *auth.Handler, controller *Controller) *chi.Mux {
	router := chi.NewRouter()
	router.Use(utils.RecoveryMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost:*", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// admin routes
	adminRouter := chi.NewRouter()
	adminRouter.Use(authHandler.AuthMiddleware)
	adminRouter.Use(authHandler.AdminCheckMiddleware)
	router.Mount("/admin", adminRouter)
	adminRouter.Post("/create-tournament", controller.createTournament)
	adminRouter.Post("/delete-tournament", controller.deleteTournament)
	adminRouter.Post("/start-tournament", controller.startTournament)

	// authenticated routes
	router.Group(func(router chi.Router) {
		router.Use(authHandler.AuthMiddleware)
		router.Post("/logout", controller.handleLogout)
		router.Get("/me", authHandler.HandleMe)
		router.Post("/set-username", controller.updateUsername)
		router.HandleFunc("/ws", controller.socketManager.ServeWS)
	})

	// routes that don't need authentication
	router.Get("/profile/{username}", controller.writeProfileInfo)
	router.Get("/games/{username}", controller.writeGames)
	router.Get("/auth/login/google", auth.GoogleLogin)
	router.Get("/auth/callback/google", authHandler.GoogleCallback)
	router.Post("/meow", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		log.Println(username)
	})
	router.Get("/get-leaderboard", controller.writeLeaderBoard)
	// game subroute

	gameRouter := chi.NewRouter()
	gameRouter.Use(authHandler.AuthMiddleware)
	gameRouter.Get("/{gameID}", controller.writeGameInfo)
	router.Mount("/game", gameRouter)

	// tournament subroute
	tournamentRouter := chi.NewRouter()
	tournamentRouter.Use(authHandler.AuthMiddleware)
	router.Mount("/tournament", tournamentRouter)
	tournamentRouter.Get("/{tournamentID}", controller.writeTournamentInfo)
	tournamentRouter.Get("/scheduled", controller.writeScheduledTournaments)
	tournamentRouter.Get("/live", controller.writeLiveTournaments)

	return router
}
