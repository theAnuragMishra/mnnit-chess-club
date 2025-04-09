package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/notnil/chess"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func (c *Controller) createGame(p1, p2 int32, p1un, p2un string, timeControl string) (*game.Game, error) {
	parts := strings.Split(timeControl, "+")
	if len(parts) != 2 {
		return nil, errors.New("not a valid time control format")
	}

	baseTime, err1 := strconv.Atoi(parts[0])
	increment, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return nil, errors.New("not a valid time control format")
	}

	createdGame := game.NewGame(time.Duration(baseTime)*time.Minute, time.Duration(increment)*time.Second, p1, p2)

	id, err := c.Queries.CreateGame(context.Background(), database.CreateGameParams{
		BaseTime:      int32(baseTime * 60),
		Increment:     int32(increment),
		WhiteID:       &p1,
		BlackID:       &p2,
		WhiteUsername: &p1un,
		BlackUsername: &p2un,
		Fen:           createdGame.Board.FEN(),
	})
	if err != nil {
		return nil, err
	}

	createdGame.ID = id
	c.GameManager.Games[id] = createdGame
	timer := time.AfterFunc(time.Second*20, func() {
		c.abortGame(createdGame)
	})

	createdGame.AbortTimer = timer
	return createdGame, nil
}

func (c *Controller) abortGame(g *game.Game) {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()

	reason := "Game Aborted"
	etl := int32(g.BaseTime.Milliseconds())

	err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
		Result:           "aborted",
		EndTimeLeftWhite: &etl,
		EndTimeLeftBlack: &etl,
		ResultReason:     &reason,
		ID:               g.ID,
	})
	if err != nil {
		log.Println("error ending game with result", err)
		return
	}
	payload, err := json.Marshal(map[string]any{"gameID": g.ID, "Result": "Aborted", "Reason": reason})
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{
		Type:    "game_abort",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.Broadcast(e)
	delete(c.GameManager.Games, g.ID)
}

func (c *Controller) handleGameTimeout(g *game.Game) {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()

	var etlw, etlb int32
	var result, reason string

	if g.Board.Position().Turn() == chess.White {
		etlw = 0
		etlb = int32(g.TimeBlack.Milliseconds())
		result = "0-1"
		reason = "White Timeout"
	} else {
		etlb = 0
		etlw = int32(g.TimeWhite.Milliseconds())
		result = "1-0"
		reason = "Black Timeout"
	}

	err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
		Result:           result,
		EndTimeLeftWhite: &etlw,
		EndTimeLeftBlack: &etlb,
		ResultReason:     &reason,
		ID:               g.ID,
	})
	if err != nil {
		log.Println("error ending game on timeout", err)
	}
	payload, err := json.Marshal(map[string]any{"Result": result, "Reason": reason, "gameID": g.ID})
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{
		Type:    "timeup",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.Broadcast(e)

	delete(c.GameManager.Games, g.ID)
}
