package control

import (
	"context"
	"errors"
	"log"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"
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
}

type Controller struct {
	SocketManager     *socket.Manager
	GameManager       *game.Manager
	Matcher           *game.Matcher
	Queries           *database.Queries
	TournamentManager *tournament.Manager
	gameRecv          chan game.EndNotification
	tournamentRecv    chan tournament.ControllerMsg
}

func NewController(queries *database.Queries) *Controller {
	c := Controller{}
	c.SocketManager = socket.NewManager(c.handleEvent, c.handleClientDisconnect)
	c.GameManager = game.NewManager()
	c.Matcher = game.NewMatcher()
	c.TournamentManager = tournament.NewManager()
	c.Queries = queries
	c.gameRecv = make(chan game.EndNotification, 256)
	c.tournamentRecv = make(chan tournament.ControllerMsg, 256)

	if err := queries.DeleteLiveGames(context.Background()); err != nil {
		log.Println("Error deleting live games", err)
	}
	if err := queries.DeleteLiveTournaments(context.Background()); err != nil {
		log.Println("Error deleting live tournaments, err")
	}
	go c.gameReceiveListener()
	go c.tournamentReceiveListener()

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

}
