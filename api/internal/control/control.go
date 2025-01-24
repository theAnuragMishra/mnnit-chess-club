package control

import (
	"errors"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

type Controller struct {
	SocketManager *socket.Manager
}

func NewController() *Controller {
	c := Controller{}
	c.SocketManager = socket.NewManager(c.HandleEvent)
	
	return &c
}
func (c *Controller) HandleEvent(event socket.Event, client *socket.Client) error {
	if handler, ok := handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no event of this type")
	}
}
