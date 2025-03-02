package control

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func Move(c *Controller, event socket.Event, client *socket.Client) error {
	// fmt.Println(string(event.Payload))
	log.Println("inside make move handler")

	var move movePayload
	if err := json.Unmarshal(event.Payload, &move); err != nil {
		log.Println("Invalid move payload")
		return err
	}
	gameID := move.GameID

	var foundGame *game.Game
	for _, g := range c.GameManager.Games {
		if g.ID == gameID {
			foundGame = g
		}
	}

	if foundGame == nil {
		return errors.New("game not found")
	}

	if foundGame.BlackID != client.UserID && foundGame.WhiteID != client.UserID {
		return errors.New("not one of the players")
	}

	if foundGame.Result != "ongoing" {
		return errors.New("game has ended")
	}

	if foundGame.WhiteID == 0 || foundGame.BlackID == 0 {
		return errors.New("game not started")
	}

	message, result := foundGame.MakeMove(client.UserID, move.MoveStr)

	if message == "game over with result" {
		// log.Println("game ho gya over")
		err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
			Result: result,
			ID:     foundGame.ID,
		})
		if err != nil {
			log.Println("error ending game with result", err)
		}
	}
	var x int32
	if foundGame.Board.Position().Turn() == 'w' {
		x = foundGame.BlackID
	} else {
		x = foundGame.WhiteID
	}
	insertedMove, err := c.Queries.InsertMove(context.Background(), database.InsertMoveParams{
		GameID:       gameID,
		MoveNumber:   int32(foundGame.GameLength + 1),
		PlayerID:     &x,
		MoveNotation: move.MoveStr,
		Orig:         move.Orig,
		Dest:         move.Dest,
		MoveFen:      foundGame.Board.FEN(),
	})
	if err != nil {
		log.Println(err)
	}
	foundGame.GameLength += 1

	err = c.Queries.UpdateGameLengthAndFEN(context.Background(), database.UpdateGameLengthAndFENParams{
		Fen:        foundGame.Board.FEN(),
		GameLength: foundGame.GameLength,
		ID:         foundGame.ID,
	})
	if err != nil {
		log.Println("error updating game fen", err)
	}

	// log.Println(foundGame.Board.Position().Turn())

	payload, err := json.Marshal(map[string]interface{}{"move": insertedMove, "Result": result, "message": message, "timeBlack": foundGame.TimeBlack.Seconds(), "timeWhite": foundGame.TimeWhite.Seconds()})
	if err != nil {
		log.Println("error marshalling new game payload")
		return nil
	}
	e := socket.Event{
		Type:    "Move_Response",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.Broadcast(e)

	return nil
}

func TimeUp(c *Controller, event socket.Event, client *socket.Client) error {
	fmt.Println("inside timeup event")
	var timeupData timeupPayload
	if err := json.Unmarshal(event.Payload, &timeupData); err != nil {
		log.Println("Invalid timeup payload")
		return err
	}

	gameID := timeupData.GameID

	var foundGame *game.Game
	for _, g := range c.GameManager.Games {
		if g.ID == gameID {
			foundGame = g
		}
	}

	if foundGame == nil {
		return errors.New("game not found")
	}

	if foundGame.Result != "ongoing" {
		return errors.New("game has ended")
	}

	moveTime := time.Since(foundGame.LastMoveTime)

	// fmt.Println(moveTime)

	if foundGame.Board.Position().Turn() == chess.White {
		if moveTime > foundGame.TimeWhite {
			foundGame.TimeWhite = 0
			foundGame.Result = "0-1"
		}
		err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
			Result: "0-1",
			ID:     foundGame.ID,
		})
		if err != nil {
			log.Println("error ending game with result", err)
		}
		payload, err := json.Marshal(map[string]interface{}{"Result": "0-1"})
		if err != nil {
			log.Println("error marshalling new game payload")
			return nil
		}
		e := socket.Event{
			Type:    "timeup",
			Payload: json.RawMessage(payload),
		}
		c.SocketManager.Broadcast(e)
	}
	if foundGame.Board.Position().Turn() == chess.Black {
		if moveTime > foundGame.TimeBlack {
			foundGame.TimeBlack = 0
			foundGame.Result = "1-0"
		}
		err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
			Result: "1-0",
			ID:     foundGame.ID,
		})
		if err != nil {
			log.Println("error ending game with result", err)
		}
		payload, err := json.Marshal(map[string]interface{}{"Result": "1-0"})
		if err != nil {
			log.Println("error marshalling new game payload")
			return nil
		}
		e := socket.Event{
			Type:    "timeup",
			Payload: json.RawMessage(payload),
		}
		c.SocketManager.Broadcast(e)

	}

	return errors.New("time not actually up")
}
