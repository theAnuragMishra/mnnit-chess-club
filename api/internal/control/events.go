package control

import "github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"

type EventHandler func(c *Controller, event socket.Event, client *socket.Client) error

var handlers = map[string]EventHandler{
	"move":             Move,
	"chat":             Chat,
	"draw":             Draw,
	"resign":           Resign,
	"init_game":        InitGame,
	"room_change":      RoomChange,
	"leave_room":       LeaveRoom,
	"create_challenge": CreateChallenge,
	"accept_challenge": AcceptChallenge,
	"join_leave":       HandleJoinLeave,
	"create_rematch":   CreateRematch,
	"accept_rematch":   AcceptRematch,
}
