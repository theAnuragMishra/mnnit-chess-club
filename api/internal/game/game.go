package game

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/notnil/chess"
)

type Result string

const (
	WhiteWins Result = "WHITE_WINS"
	BlackWins Result = "BLACK_WINS"
	Draw      Result = "DRAW"
)

type Game struct {
	Id        string
	Player1Id string
	Player2Id string
	Board     *chess.Game
	moveCount int
	Result    Result
}

func NewGame(player1 string, player2 string) *Game {
	var board *chess.Game
	return &Game{
		Id:        uuid.New().String(),
		Player1Id: player1,
		Player2Id: player2,
		Board:     board,
	}
}

func (g *Game) MakeMove(player string, move string) {
	if g.Board.Position().Turn() == chess.White && player != g.Player1Id {
		return
	}
	if g.Board.Position().Turn() == chess.Black && player != g.Player2Id {
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

	fmt.Println(g.Board.Position())
	fmt.Println(g.Board.Position().Board())
	fmt.Println(g.Board)

	// moveStr := "{\"move\": " + move + "}"

	//g.Socket.Broadcast(socket.Event{
	//	Type:    socket.EventMove,
	//	Payload: json.RawMessage(moveStr),
	//})
}
