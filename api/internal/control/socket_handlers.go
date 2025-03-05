package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func Move(c *Controller, event socket.Event, client *socket.Client) error {
	// fmt.Println(string(event.Payload))

	var move movePayload
	if err := json.Unmarshal(event.Payload, &move); err != nil {
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
		etlb := int32(foundGame.TimeBlack.Seconds())
		etlw := int32(foundGame.TimeWhite.Seconds())

		err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
			Result:           result,
			ID:               foundGame.ID,
			EndTimeLeftBlack: &etlb,
			EndTimeLeftWhite: &etlw,
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
	var timeupData timeupPayload
	if err := json.Unmarshal(event.Payload, &timeupData); err != nil {
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
		etlb := int32(foundGame.TimeBlack.Seconds())
		etlw := int32(foundGame.TimeWhite.Seconds())
		err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
			Result:           "0-1",
			ID:               foundGame.ID,
			EndTimeLeftBlack: &etlb,
			EndTimeLeftWhite: &etlw,
		})
		if err != nil {
			log.Println("error ending game with result", err)
		}
		payload, err := json.Marshal(map[string]interface{}{"Result": "0-1"})
		if err != nil {
			return err
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
		etlb := int32(foundGame.TimeBlack.Seconds())
		etlw := int32(foundGame.TimeWhite.Seconds())
		err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
			Result:           "1-0",
			ID:               foundGame.ID,
			EndTimeLeftWhite: &etlw,
			EndTimeLeftBlack: &etlb,
		})
		if err != nil {
			log.Println("error ending game with result", err)
		}
		payload, err := json.Marshal(map[string]interface{}{"Result": "1-0"})
		if err != nil {
			return err
		}
		e := socket.Event{
			Type:    "timeup",
			Payload: json.RawMessage(payload),
		}
		c.SocketManager.Broadcast(e)

	}

	return errors.New("time not actually up")
}

func Chat(c *Controller, event socket.Event, client *socket.Client) error {
	var chat ChatPayload
	if err := json.Unmarshal(event.Payload, &chat); err != nil {
		return errors.New("invalid payload")
	}

	if chat.Sender == chat.Receiver {
		return errors.New("same Sender and Receiver")
	}

	if chat.Sender != client.UserID {
		return errors.New("client not the sender")
	}
	var foundGame *game.Game
	for _, g := range c.GameManager.Games {
		if g.ID == chat.GameID {
			foundGame = g
		}
	}

	if foundGame == nil {
		return errors.New("game not found")
	}
	// fmt.Println(foundGame.Result)
	if foundGame.Result != "ongoing" {
		return errors.New("game has ended")
	}

	if (chat.Sender != foundGame.WhiteID && chat.Sender != foundGame.BlackID) || (chat.Receiver != foundGame.WhiteID && chat.Receiver != foundGame.BlackID) {
		return errors.New("game id doesn't correspond")
	}

	otherClient := c.SocketManager.FindClientByUserID(chat.Receiver)

	payload, err := json.Marshal(map[string]interface{}{"sender": chat.SenderUsername, "receiver": chat.ReceiverUsername, "gameID": chat.GameID, "text": chat.Text})
	if err != nil {
		return errors.New("error marshalling chat payload")
	}

	e := socket.Event{
		Type:    "chat",
		Payload: json.RawMessage(payload),
	}

	client.Send(e)
	otherClient.Send(e)

	return nil
}
