package control

import "github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"

type EventHandler func(c *Controller, event socket.Event, client *socket.Client) error

var handlers = map[string]EventHandler{
	"init_game": InitGame,
	"move":      Move,
}
