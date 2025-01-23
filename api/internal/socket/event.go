package socket

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage          = "send_message"
	EventInitGame             = "init_game"
	EventMove                 = "move"
	EventGameOver             = "game_over"
	EventJoinGame             = "join_game"
	EventOpponentDisconnected = "opponent_disconnected"
	EventJoinRoom             = "join_room"
	EventGameNotFound         = "game_not_found"
	EventGameJoined           = "game_joined"
	EventGameEnded            = "game_ended"
	EventGameAlert            = "game_alert"
	EventGameAdded            = "game_added"
	EventGameTime             = "game_time"
	EventExitGame             = "exit_game"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type MoveEvent struct {
}
