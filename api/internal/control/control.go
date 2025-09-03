package control

import (
	"context"
	"errors"
	"log"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

type EventHandler func(c *Controller, event socket.Event, client *socket.Client) error

var handlers = map[string]EventHandler{
	"move":             move,
	"chat":             chat,
	"draw":             draw,
	"resign":           resign,
	"init_game":        initGame,
	"room_change":      roomChange,
	"leave_room":       leaveRoom,
	"create_challenge": createChallenge,
	"accept_challenge": acceptChallenge,
	"join_leave":       handleJoinLeave,
	"rematch":          rematch,
	"berserk":          berserk,
}

type Controller struct {
	socketManager     *socket.Manager
	gameManager       *gameManager
	matcher           *matcher
	rematchManager    *rematchManager
	challengeManager  *challengeManager
	queries           *database.Queries
	tournamentManager *tournamentManager
}

func NewController(queries *database.Queries) *Controller {
	c := Controller{}
	c.socketManager = socket.NewManager(c.handleEvent, c.handleClientDisconnect)
	c.gameManager = newGameManager()
	c.matcher = newMatcher()
	c.rematchManager = newRematchManager()
	c.challengeManager = newChallengeManager()
	c.tournamentManager = newTournamentManager()
	c.queries = queries

	if err := queries.DeleteLiveGames(context.Background()); err != nil {
		log.Println("Error deleting live games", err)
	}
	if err := queries.DeleteLiveTournaments(context.Background()); err != nil {
		log.Println("Error deleting live tournaments, err")
	}
	return &c
}

func (c *Controller) handleEvent(event socket.Event, client *socket.Client) error {
	if handler, ok := handlers[event.Type]; ok {
		if err := handler(c, event, client); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no event of this type")
	}
}

func (c *Controller) handleClientDisconnect(client *socket.Client) {
	c.matcher.Lock()
	defer c.matcher.Unlock()
	for i := range 12 {
		c.matcher.removeClient(client, i)
	}
}
