package game

import (
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"log"

	"github.com/notnil/chess"
)

type Game struct {
	ID      int32
	WhiteID int32
	BlackID int32
	Board   *chess.Game
}

func NewGame(player1 int32, player2 int32) *Game {
	board := chess.NewGame()
	return &Game{

		WhiteID: player1,
		BlackID: player2,
		Board:   board,
	}
}

func DatabaseGameToGame(game *database.Game) *Game {
	fen, _ := chess.FEN(game.Fen)
	return &Game{
		ID:      game.ID,
		WhiteID: *game.WhiteID,
		BlackID: *game.BlackID,
		Board:   chess.NewGame(fen),
	}
}

func (g *Game) MakeMove(player int32, move string) (string, string) {
	if g.Board.Position().Turn() == chess.White && player != g.WhiteID {
		return "not your turn", ""
	}
	if g.Board.Position().Turn() == chess.Black && player != g.BlackID {
		return "not your turn", ""
	}

	if g.Board.Outcome() != "*" {

		return "game has ended", string(g.Board.Outcome())
	}

	if err := g.Board.MoveStr(move); err != nil {
		log.Println(err)
		return "error making move", ""
	}
	log.Println(g.Board.Position().Board().Draw())

	if g.Board.Outcome() != "*" {
		return "game over with result", string(g.Board.Outcome())
	}

	// fmt.Println(g.Board.Position())
	// fmt.Println(g.Board.Position().Board())
	// fmt.Println(g.Board)

	return "move successful", ""

	// moveStr := "{\"move\": " + move + "}"

	//g.Socket.Broadcast(socket.Event{
	//	Type:    socket.EventMove,
	//	Payload: json.RawMessage(moveStr),
	//})
}
