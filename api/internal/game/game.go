package game

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/notnil/chess"
)

type Result string

type Game struct {
	ID                  int32
	WhitePlayerID       uuid.UUID
	BlackPlayerID       uuid.UUID
	WhitePlayerUsername string
	BlackPlayerUsername string
	Board               *chess.Game
	moveCount           int
	Result              Result
}

func NewGame(id int32, player1ID uuid.UUID, player1Username string, player2ID uuid.UUID, player2Username string) *Game {
	board := chess.NewGame()
	return &Game{
		ID:                  id,
		WhitePlayerID:       player1ID,
		WhitePlayerUsername: player1Username,
		BlackPlayerID:       player2ID,
		BlackPlayerUsername: player2Username,
		Board:               board,
	}
}

func (g *Game) MakeMove(playerID uuid.UUID, move string) Result {
	if g.Board.Position().Turn() == chess.White && playerID != g.WhitePlayerID {
		return "not your turn"
	}
	if g.Board.Position().Turn() == chess.Black && playerID != g.BlackPlayerID {
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
	fmt.Println(g.Board.Position().Board().Draw())

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
