package control

import (
	"fmt"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func InitGameHandler(event socket.Event, client *socket.Client) error {
	fmt.Println(event)
	return nil
}
