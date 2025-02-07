package control

import (
	"errors"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

type Controller struct {
	SocketManager *socket.Manager
	GameManager   *game.Manager
}

func NewController(authHandler *auth.Handler) *Controller {
	c := Controller{}
	c.SocketManager = socket.NewManager(c.HandleEvent, authHandler)
	c.GameManager = game.NewManager()

	return &c
}

func (c *Controller) HandleEvent(event socket.Event, client *socket.Client) error {
	if handler, ok := handlers[event.Type]; ok {
		if err := handler(c, event, client); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no event of this type")
	}
}
