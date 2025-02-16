package game

import (
	"log"

	"github.com/notnil/chess"
)

type Result string

type Game struct {
	ID        int32
	WhiteID   int32
	BlackID   int32
	Board     *chess.Game
	moveCount int
	Result    Result
}

func NewGame(player1 int32, player2 int32) *Game {
	board := chess.NewGame()
	return &Game{

		WhiteID: player1,
		BlackID: player2,
		Board:   board,
	}
}

func (g *Game) MakeMove(player int32, move string) Result {
	if g.Board.Position().Turn() == chess.White && player != g.WhiteID {
		return "not your turn"
	}
	if g.Board.Position().Turn() == chess.Black && player != g.BlackID {
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
	log.Println(g.Board.Position().Board().Draw())

	if g.Board.Outcome() != "*" {
		g.Result = Result(g.Board.Outcome())
		return Result(g.Board.Outcome())
	}

	// fmt.Println(g.Board.Position())
	// fmt.Println(g.Board.Position().Board())
	// fmt.Println(g.Board)

	return "move successful"

	// moveStr := "{\"move\": " + move + "}"

	//g.Socket.Broadcast(socket.Event{
	//	Type:    socket.EventMove,
	//	Payload: json.RawMessage(moveStr),
	//})
}
