package game

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/notnil/chess"
)

type Result string

type Game struct {
	ID        string
	Player1Id string
	Player2Id string
	Board     *chess.Game
	moveCount int
	Result    Result
}

func NewGame(player1 string) *Game {
	board := chess.NewGame()
	return &Game{
		ID:        uuid.New().String(),
		Player1Id: player1,
		Board:     board,
	}
}

func (g *Game) MakeMove(player string, move string) Result {
	if g.Board.Position().Turn() == chess.White && player != g.Player1Id {
		return "not your turn"
	}
	if g.Board.Position().Turn() == chess.Black && player != g.Player2Id {
		return "not your turn"
	}

	if g.Result != "" {
		log.Println("Trying to make a move after game has finished")
		return g.Result
	}
	if err := g.Board.MoveStr(move); err != nil {
		log.Println(err)
		return "error making move"
	}

	if g.Board.Outcome() != "*" {
		g.Result = Result(g.Board.Outcome())
		return Result(g.Board.Outcome())
	}

	// fmt.Println(g.Board.Position())
	// fmt.Println(g.Board.Position().Board())
	// fmt.Println(g.Board)
	fmt.Println(g.Board.Position().Board().Draw())

	return "move successful"

	// moveStr := "{\"move\": " + move + "}"

	//g.Socket.Broadcast(socket.Event{
	//	Type:    socket.EventMove,
	//	Payload: json.RawMessage(moveStr),
	//})
}
