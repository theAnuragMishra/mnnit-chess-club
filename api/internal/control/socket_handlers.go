package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func InitGame(c *Controller, event socket.Event, client *socket.Client) error {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()
	var initGamePayload InitGamePayload
	if err := json.Unmarshal(event.Payload, &initGamePayload); err != nil {
		return err
	}
	pendingUser, exists := c.GameManager.PendingUsers[initGamePayload.TimeControl]
	if !exists {
		log.Println("no pending game, creating...")
		c.GameManager.PendingUsers[initGamePayload.TimeControl] = client.UserID
	} else {
		log.Println("match found...")
		delete(c.GameManager.PendingUsers, initGamePayload.TimeControl)
		if pendingUser == client.UserID {
			return errors.New("same player tryna play both sides")
		}
		u1, err1 := c.Queries.GetUsernameAndRating(context.Background(), pendingUser)
		u2, err2 := c.Queries.GetUsernameAndRating(context.Background(), client.UserID)
		if err1 != nil || err2 != nil {
			return errors.New("server error while fetching usernames")
		}
		createdGame, err := c.createGame(pendingUser, client.UserID, *u1.Username, *u2.Username, initGamePayload.TimeControl, u1.Rating, u2.Rating)
		if err != nil {
			return err
		}
		payload := map[string]any{"GameID": createdGame.ID}
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		e := socket.Event{Type: "GoToGame", Payload: json.RawMessage(rawPayload)}
		otherClient := c.SocketManager.FindClientByUserID(pendingUser)
		if client != nil {
			client.Send(e)
		}
		if otherClient != nil {
			otherClient.Send(e)
		}
	}
	return nil
}

func Move(c *Controller, event socket.Event, client *socket.Client) error {
	// fmt.Println(string(event.Payload))

	var move movePayload
	if err := json.Unmarshal(event.Payload, &move); err != nil {
		return err
	}
	gameID := move.GameID

	foundGame, exists := c.GameManager.Games[gameID]

	if !exists {
		return errors.New("game not found")
	}

	if foundGame.Result != "ongoing" {
		return errors.New("game has ended")
	}

	if foundGame.Board.Position().Turn() == chess.White && client.UserID != foundGame.WhiteID {
		return errors.New("not your turn")
	}
	if foundGame.Board.Position().Turn() == chess.Black && client.UserID != foundGame.BlackID {
		return errors.New("not your turn")
	}

	message, result := foundGame.MakeMove(move.MoveStr)

	if message == "error making move" {
		return errors.New("error making move")
	}

	foundGame.GameLength += 1

	if foundGame.AbortTimer != nil {
		if foundGame.GameLength == 1 {
			foundGame.AbortTimer.Reset(time.Second * 20)
		} else {
			foundGame.AbortTimer.Stop()
			foundGame.AbortTimer = nil
		}
	}

	if foundGame.GameLength == 2 {
		var t time.Duration
		if foundGame.Board.Position().Turn() == chess.White {
			t = foundGame.TimeWhite
		} else {
			t = foundGame.TimeBlack
		}
		timer := time.AfterFunc(t, func() { c.handleGameTimeout(foundGame) })
		foundGame.ClockTimer = timer
	} else if foundGame.GameLength > 2 {
		if foundGame.Board.Position().Turn() == chess.White {
			foundGame.ClockTimer.Reset(foundGame.TimeWhite)
		} else {
			foundGame.ClockTimer.Reset(foundGame.TimeBlack)
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
		MoveNumber:   int32(foundGame.GameLength),
		PlayerID:     &x,
		MoveNotation: move.MoveStr,
		Orig:         move.Orig,
		Dest:         move.Dest,
		MoveFen:      foundGame.Board.FEN(),
	})
	if err != nil {
		log.Println(err)
	}

	err = c.Queries.UpdateGameLengthAndFEN(context.Background(), database.UpdateGameLengthAndFENParams{
		Fen:        foundGame.Board.FEN(),
		GameLength: foundGame.GameLength,
		ID:         foundGame.ID,
	})
	if err != nil {
		log.Println("error updating game fen", err)
	}

	// log.Println(foundGame.Board.Position().Turn())
	var cw, cb int
	if result != "" {
		// log.Println("game ho gya over")

		etlb := int32(foundGame.TimeBlack.Milliseconds())
		etlw := int32(foundGame.TimeWhite.Milliseconds())

		cw, cb, err = c.endGame(foundGame.ID, &etlw, &etlb, result, &message, foundGame.WhiteID, foundGame.BlackID)
		if err != nil {
			log.Println("error ending game with result", err)
		}
		foundGame.ClockTimer.Stop()
		delete(c.GameManager.Games, gameID)
	}
	payload, err := json.Marshal(map[string]any{"gameID": gameID, "move": insertedMove, "Result": result, "message": message, "timeBlack": foundGame.TimeBlack.Milliseconds(), "timeWhite": foundGame.TimeWhite.Milliseconds(), "changeW": cw, "changeB": cb})
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

	foundGame, exists := c.GameManager.Games[chat.GameID]
	if !exists {
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

	payload, err := json.Marshal(map[string]any{"sender": chat.SenderUsername, "receiver": chat.ReceiverUsername, "gameID": chat.GameID, "text": chat.Text})
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

func Draw(c *Controller, event socket.Event, client *socket.Client) error {
	var draw DRPayload
	if err := json.Unmarshal(event.Payload, &draw); err != nil {
		return err
	}

	if client.UserID != draw.PlayerID {
		return errors.New("not the player")
	}

	gameID := draw.GameID

	foundGame, exists := c.GameManager.Games[gameID]

	if !exists {
		return errors.New("game not found")
	}

	if foundGame.GameLength < 2 {
		return errors.New("cannot resign a game where one or both sides haven't played")
	}

	if foundGame.Result != "ongoing" {
		return errors.New("game has ended")
	}

	if foundGame.WhiteID != draw.PlayerID && foundGame.BlackID != draw.PlayerID {
		return errors.New("not one of the players")
	}

	if foundGame.DrawOfferedBy == 0 {
		foundGame.DrawOfferedBy = draw.PlayerID
		payload, err := json.Marshal(map[string]any{"gameID": draw.GameID})
		if err != nil {
			return errors.New("error marshalling chat payload")
		}

		e := socket.Event{
			Type:    "drawOffer",
			Payload: json.RawMessage(payload),
		}

		var other int32
		if draw.PlayerID == foundGame.BlackID {
			other = foundGame.WhiteID
		} else {
			other = foundGame.BlackID
		}

		otherClient := c.SocketManager.FindClientByUserID(other)

		otherClient.Send(e)
	} else if foundGame.DrawOfferedBy != draw.PlayerID {
		reason := "Draw by mutual agreement"
		timeTaken := time.Since(foundGame.LastMoveTime)

		if foundGame.Board.Position().Turn() == chess.White {
			foundGame.TimeWhite -= timeTaken
		} else {
			foundGame.TimeBlack -= timeTaken
		}

		etlb := int32(foundGame.TimeBlack.Milliseconds())
		etlw := int32(foundGame.TimeWhite.Milliseconds())

		cw, cb, err := c.endGame(foundGame.ID, &etlw, &etlb, "1/2-1/2", &reason, foundGame.WhiteID, foundGame.BlackID)
		if err != nil {
			log.Println("error ending game with result", err)
		}

		payload, err := json.Marshal(map[string]any{"gameID": gameID, "Result": "1/2-1/2", "Reason": reason, "changeW": cw, "changeB": cb})
		if err != nil {
			return err
		}
		e := socket.Event{
			Type:    "gameDrawn",
			Payload: json.RawMessage(payload),
		}
		c.SocketManager.Broadcast(e)
		foundGame.ClockTimer.Stop()
		delete(c.GameManager.Games, gameID)
	}

	return nil
}

func Resign(c *Controller, event socket.Event, client *socket.Client) error {
	var resign DRPayload
	if err := json.Unmarshal(event.Payload, &resign); err != nil {
		return err
	}

	if client.UserID != resign.PlayerID {
		return errors.New("not the player")
	}

	gameID := resign.GameID

	foundGame, exists := c.GameManager.Games[gameID]

	if !exists {
		return errors.New("game not found")
	}

	if foundGame.GameLength < 2 {
		return errors.New("cannot resign a game where one or both sides haven't played")
	}

	if foundGame.Result != "ongoing" {
		return errors.New("game has ended")
	}

	if foundGame.WhiteID != resign.PlayerID && foundGame.BlackID != resign.PlayerID {
		return errors.New("not one of the players")
	}

	var result string
	var reason string

	if foundGame.WhiteID == resign.PlayerID {
		result = "0-1"
		reason = "White Resigned"
	} else {
		result = "1-0"
		reason = "Black Resigned"
	}

	timeTaken := time.Since(foundGame.LastMoveTime)

	if foundGame.Board.Position().Turn() == chess.White {
		foundGame.TimeWhite -= timeTaken
	} else {
		foundGame.TimeBlack -= timeTaken
	}

	etlb := int32(foundGame.TimeBlack.Milliseconds())
	etlw := int32(foundGame.TimeWhite.Milliseconds())

	cw, cb, err := c.endGame(foundGame.ID, &etlw, &etlb, result, &reason, foundGame.WhiteID, foundGame.BlackID)
	if err != nil {
		log.Println("error ending game with result", err)
	}

	payload, err := json.Marshal(map[string]any{"gameID": gameID, "Result": result, "Reason": reason, "changeW": cw, "changeB": cb})
	if err != nil {
		return err
	}
	e := socket.Event{
		Type:    "resignation",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.Broadcast(e)

	foundGame.ClockTimer.Stop()
	delete(c.GameManager.Games, gameID)
	return nil
}
