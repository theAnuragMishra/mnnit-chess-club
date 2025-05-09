package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"

	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func RoomChange(c *Controller, event socket.Event, client *socket.Client) error {
	c.SocketManager.Lock()
	defer c.SocketManager.Unlock()
	var payload RoomPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}
	if client.Room == payload.RoomID {
		return nil
	}
	delete(c.SocketManager.Rooms[client.Room], client)
	client.Room = payload.RoomID
	if c.SocketManager.Rooms[payload.RoomID] == nil {
		c.SocketManager.Rooms[payload.RoomID] = make(map[*socket.Client]bool)
	}
	c.SocketManager.Rooms[payload.RoomID][client] = true
	// printing clients for debugging purposes
	// for _, x := range c.SocketManager.Rooms {
	// 	for y := range x {
	// 		fmt.Println(y)
	// 	}
	// }
	return nil
}

func LeaveRoom(c *Controller, event socket.Event, client *socket.Client) error {
	c.SocketManager.Lock()
	defer c.SocketManager.Unlock()
	var payload RoomPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}
	delete(c.SocketManager.Rooms[client.Room], client)
	client.Room = ""
	// for _, x := range c.SocketManager.Rooms {
	// 	for y := range x {
	// 		fmt.Println(y)
	// 	}
	// }
	return nil
}

func InitGame(c *Controller, event socket.Event, client *socket.Client) error {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()
	var timeControl game.TimeControl
	if err := json.Unmarshal(event.Payload, &timeControl); err != nil {
		return err
	}
	//log.Println(timeControl)
	if timeControl.BaseTime <= 0 || timeControl.BaseTime > 10800 || timeControl.Increment < 0 || timeControl.Increment > 180 {
		return errors.New("invalid time control")
	}
	pendingUser, exists := c.GameManager.PendingUsers[timeControl]
	if !exists {
		log.Println("no pending game, creating...")
		c.GameManager.PendingUsers[timeControl] = client.UserID
	} else {
		log.Println("game found...")
		delete(c.GameManager.PendingUsers, timeControl)
		if pendingUser == client.UserID {
			return nil
		}
		rating1, err1 := c.Queries.GetUserRating(context.Background(), pendingUser)
		rating2, err2 := c.Queries.GetUserRating(context.Background(), client.UserID)
		if err1 != nil || err2 != nil {
			return errors.New("server error while fetching ratings")
		}
		id, err := c.generateUniqueGameID()
		if err != nil {
			return err
		}
		createdGame, err := c.createGame(id, pendingUser, client.UserID, timeControl, rating1, rating2)
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

		client.Send(e)

		if otherClient != nil {
			otherClient.Send(e)
		}
	}
	return nil
}

func CreateChallenge(c *Controller, event socket.Event, client *socket.Client) error {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()
	var timeControl game.TimeControl
	if err := json.Unmarshal(event.Payload, &timeControl); err != nil {
		return err
	}
	//log.Println(timeControl)
	if timeControl.BaseTime <= 0 || timeControl.BaseTime > 10800 || timeControl.Increment < 0 || timeControl.Increment > 180 {
		return errors.New("invalid time control")
	}
	id, err := c.generateUniqueGameID()
	if err != nil {
		return err
	}
	c.GameManager.PendingChallenges[id] = game.Challenge{
		TimeControl:     timeControl,
		Creator:         client.UserID,
		CreatorUsername: client.Username,
	}
	payload := map[string]any{"GameID": id}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "GoToGame", Payload: json.RawMessage(rawPayload)}
	client.Send(e)
	return nil
}

func AcceptChallenge(c *Controller, event socket.Event, client *socket.Client) error {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()
	var acceptChallengePayload AcceptChallengePayload
	if err := json.Unmarshal(event.Payload, &acceptChallengePayload); err != nil {
		return err
	}
	challenge, exists := c.GameManager.PendingChallenges[acceptChallengePayload.GameID]
	if !exists {
		return errors.New("challenge not found")
	}
	if client.UserID == challenge.Creator {
		return errors.New("same player tryna play both sides")
	}
	rating1, err1 := c.Queries.GetUserRating(context.Background(), challenge.Creator)
	rating2, err2 := c.Queries.GetUserRating(context.Background(), client.UserID)
	if err1 != nil || err2 != nil {
		return errors.New("server error while fetching ratings")
	}
	createdGame, err := c.createGame(acceptChallengePayload.GameID, challenge.Creator, client.UserID, challenge.TimeControl, rating1, rating2)
	if err != nil {
		return err
	}
	payload := map[string]any{"GameID": createdGame.ID}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "RefreshGame", Payload: json.RawMessage(rawPayload)}
	c.SocketManager.BroadcastToRoom(e, acceptChallengePayload.GameID)
	delete(c.GameManager.PendingChallenges, acceptChallengePayload.GameID)
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
		if foundGame.Board.Position().Turn() == chess.White {
			timer := time.AfterFunc(foundGame.TimeWhite, func() { c.handleGameTimeout(foundGame) })
			foundGame.ClockTimer = timer
		} else {
			timer := time.AfterFunc(foundGame.TimeBlack, func() { c.handleGameTimeout(foundGame) })
			foundGame.ClockTimer = timer
		}
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
	c.SocketManager.BroadcastToRoom(e, gameID)

	return nil
}

func Chat(c *Controller, event socket.Event, client *socket.Client) error {
	var chat ChatPayload
	if err := json.Unmarshal(event.Payload, &chat); err != nil {
		return errors.New("invalid payload")
	}

	payload, err := json.Marshal(map[string]any{"sender": client.Username, "gameID": client.Room, "text": chat.Text})
	if err != nil {
		return errors.New("error marshalling chat payload")
	}

	e := socket.Event{
		Type:    "chat",
		Payload: json.RawMessage(payload),
	}

	foundGame, exists := c.GameManager.Games[client.Room]
	if !exists || foundGame.Result != "ongoing" {
		// handle game ended message
		c.SocketManager.BroadcastToRoom(e, client.Room)
		return nil
	}

	whiteClient := c.SocketManager.FindClientByUserID(foundGame.WhiteID)
	blackClient := c.SocketManager.FindClientByUserID(foundGame.BlackID)

	if client.UserID != foundGame.WhiteID && client.UserID != foundGame.BlackID {
		// handle message by non player
		c.SocketManager.BroadcastToNonPlayers(e, client.Room, whiteClient, blackClient)
		return nil
	}

	whiteClient.Send(e)
	blackClient.Send(e)
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
		c.SocketManager.BroadcastToRoom(e, gameID)
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
	c.SocketManager.BroadcastToRoom(e, gameID)

	foundGame.ClockTimer.Stop()
	delete(c.GameManager.Games, gameID)
	return nil
}
