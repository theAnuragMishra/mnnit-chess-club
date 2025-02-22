package game

import (
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"log"
	"time"

	"github.com/notnil/chess"
)

type Game struct {
	ID           int32
	Result       string
	BaseTime     time.Duration
	Increment    time.Duration
	WhiteID      int32
	BlackID      int32
	Board        *chess.Game
	GameLength   int16
	LastMoveTime time.Time
	TimeWhite    time.Duration
	TimeBlack    time.Duration
}

func NewGame(baseTime time.Duration, increment time.Duration, player1 int32, player2 int32) *Game {
	board := chess.NewGame()

	return &Game{
		Result:       "ongoing",
		TimeWhite:    baseTime,
		TimeBlack:    baseTime,
		BaseTime:     baseTime,
		Increment:    increment,
		WhiteID:      player1,
		BlackID:      player2,
		Board:        board,
		LastMoveTime: time.Now(),
	}
}

func DatabaseGameToGame(game *database.Game) *Game {
	fen, _ := chess.FEN(game.Fen)
	return &Game{
		ID:           game.ID,
		Result:       game.Result,
		BaseTime:     time.Duration(game.BaseTime) * time.Second,
		Increment:    time.Duration(game.Increment) * time.Second,
		TimeBlack:    time.Duration(game.BaseTime) * time.Second,
		TimeWhite:    time.Duration(game.BaseTime) * time.Second,
		WhiteID:      *game.WhiteID,
		BlackID:      *game.BlackID,
		Board:        chess.NewGame(fen),
		GameLength:   game.GameLength,
		LastMoveTime: time.Now(),
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

	moveTime := time.Since(g.LastMoveTime)

	//fmt.Println(moveTime)

	if g.Board.Position().Turn() == chess.White {
		if moveTime > g.TimeWhite {
			g.TimeWhite = 0
			g.Result = "0-1"
			return "game over with result", "0-1"
		} else {
			g.TimeWhite -= moveTime
			g.TimeWhite += g.Increment
		}

	}
	if g.Board.Position().Turn() == chess.Black {
		if moveTime > g.TimeBlack {
			g.TimeBlack = 0
			g.Result = "0-1"
			return "game over with result", "1-0"
		} else {
			g.TimeBlack -= moveTime
			g.TimeBlack += g.Increment
		}
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
	g.LastMoveTime = time.Now()

	return "move successful", ""

	// moveStr := "{\"move\": " + move + "}"

	//g.Socket.Broadcast(socket.Event{
	//	Type:    socket.EventMove,
	//	Payload: json.RawMessage(moveStr),
	//})
}
