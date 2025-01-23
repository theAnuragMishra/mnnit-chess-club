package game

import (
	"encoding/json"
	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"log"
)

type Result string

const (
	WhiteWins Result = "WHITE_WINS"
	BlackWins Result = "BLACK_WINS"
	Draw      Result = "DRAW"
)

type Game struct {
	GameId    string
	Player1   *socket.Client
	Player2   *socket.Client
	Board     *chess.Game
	moveCount int
	Result    Result
	Socket    *socket.Manager
}

func NewGame(gameId string, player1 *socket.Client, player2 *socket.Client, manager *socket.Manager) *Game {
	var board *chess.Game
	return &Game{
		GameId:  gameId,
		Player1: player1,
		Player2: player2,
		Board:   board,
		Socket:  manager,
	}
}

func (g *Game) MakeMove(client *socket.Client, move string) {
	if g.Board.Position().Turn() == chess.White && client != g.Player1 {
		return
	}
	if g.Board.Position().Turn() == chess.Black && client != g.Player2 {
		return
	}

	if g.Board.Outcome() != "*" {
		log.Println("Trying to make a move after game has finished")
		return
	}
	if err := g.Board.MoveStr(move); err != nil {
		log.Println(err)
		return
	}

	moveStr := "{\"move\": " + move + "}"

	g.Socket.Broadcast(socket.Event{
		Type:    socket.EventMove,
		Payload: json.RawMessage(moveStr),
	})
}
