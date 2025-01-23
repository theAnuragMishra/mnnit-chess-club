package game

import (
	"github.com/notnil/chess"
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
	Player1Id string
	Player2Id string
	Board     *chess.Game
	moveCount int
	Result    Result
}

func NewGame(gameId string, player1 string, player2 string) *Game {
	var board *chess.Game
	return &Game{
		GameId:    gameId,
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

	//moveStr := "{\"move\": " + move + "}"

	//g.Socket.Broadcast(socket.Event{
	//	Type:    socket.EventMove,
	//	Payload: json.RawMessage(moveStr),
	//})
}
