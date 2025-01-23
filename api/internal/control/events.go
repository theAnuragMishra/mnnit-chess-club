package control

import "github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"

type EventHandler func(event socket.Event, c *socket.Client) error

var handlers = map[string]EventHandler{
	"init_game": InitGameHandler,
}
