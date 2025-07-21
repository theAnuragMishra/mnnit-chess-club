package control

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func (c *Controller) WriteGameInfo(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "gameID")

	// fmt.Println(gameID)
	foundGame, err := c.Queries.GetGameInfo(r.Context(), gameID)
	if err != nil {
		if challenge, exists := c.GameManager.PendingChallenges[gameID]; exists {
			utils.RespondWithJSON(w, http.StatusOK, challenge)
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}
	serverGame, exists := c.GameManager.Games[gameID]
	if !exists {
		//database game response
		moves, err := c.Queries.GetGameMoves(r.Context(), gameID)
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError, "error getting game moves")
			return
		}
		var timeBlack, timeWhite int32
		if foundGame.EndTimeLeftWhite != nil {
			timeWhite = *foundGame.EndTimeLeftWhite
		}
		if foundGame.EndTimeLeftBlack != nil {
			timeBlack = *foundGame.EndTimeLeftBlack
		}
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"moves": moves, "game": foundGame, "timeWhite": timeWhite, "timeBlack": timeBlack})

	} else {
		//server game response
		timePassed := time.Since(serverGame.LastMoveTime)
		if serverGame.Board.Position().Turn() == chess.White {
			serverGame.TimeWhite = max(serverGame.TimeWhite-timePassed, 0)
		} else {
			serverGame.TimeBlack = max(serverGame.TimeBlack-timePassed, 0)
		}
		serverGame.LastMoveTime = time.Now()
		timeWhite := int32(serverGame.TimeWhite.Milliseconds())
		timeBlack := int32(serverGame.TimeBlack.Milliseconds())
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"moves": serverGame.Moves, "game": foundGame, "timeWhite": timeWhite, "timeBlack": timeBlack})
	}
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
	p, exists := c.GameManager.PendingUsers[timeControl]
	if !exists {
		//log.Println("no pending game, creating...")
		c.GameManager.PendingUsers[timeControl] = game.PendingUser{
			ID:     client.UserID,
			Client: client,
		}
	} else {
		//log.Println("game found...")
		delete(c.GameManager.PendingUsers, timeControl)
		if p.ID == client.UserID {
			return nil
		}
		rating1, err1 := c.Queries.GetUserRating(context.Background(), p.ID)
		rating2, err2 := c.Queries.GetUserRating(context.Background(), client.UserID)
		if err1 != nil || err2 != nil {
			return errors.New("server error while fetching ratings")
		}
		id, err := c.generateUniqueGameID()
		if err != nil {
			return err
		}
		createdGame, err := c.createGame(id, p.ID, client.UserID, timeControl, rating1, rating2, "")
		if err != nil {
			return err
		}
		payload := map[string]any{"ID": createdGame.ID, "Type": "game"}
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}
		client.Send(e)
		p.Client.SendIfConnected(e)
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
	payload := map[string]any{"ID": id, "Type": "game"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}
	client.Send(e)
	return nil
}

func AcceptChallenge(c *Controller, event socket.Event, client *socket.Client) error {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()
	var acceptChallengePayload GameIDPayload
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
	createdGame, err := c.createGame(acceptChallengePayload.GameID, challenge.Creator, client.UserID, challenge.TimeControl, rating1, rating2, "")
	if err != nil {
		return err
	}
	payload := map[string]any{"ID": createdGame.ID, "Type": "game"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
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

	var timeLeft int32
	if foundGame.Board.Position().Turn() == chess.White {
		timeLeft = int32(foundGame.TimeWhite.Milliseconds())
	} else {
		timeLeft = int32(foundGame.TimeBlack.Milliseconds())
	}

	message, result := foundGame.MakeMove(move.MoveStr)

	if message == "error making move" {
		return errors.New("error making move")
	}

	moveToSend := game.Move{
		MoveNotation: move.MoveStr,
		Orig:         move.Orig,
		Dest:         move.Dest,
		MoveFen:      foundGame.Board.FEN(),
		TimeLeft:     &timeLeft,
	}

	foundGame.Moves = append(foundGame.Moves, moveToSend)

	if foundGame.AbortTimer != nil {
		if len(foundGame.Moves) == 1 {
			if foundGame.BaseTime >= time.Second*20 {
				foundGame.AbortTimer.Reset(time.Second * 20)
			} else if foundGame.BaseTime >= time.Second*10 {
				foundGame.AbortTimer.Reset(time.Second * 10)
			} else {
				foundGame.AbortTimer.Reset(foundGame.BaseTime)
			}
		} else {
			foundGame.AbortTimer.Stop()
			foundGame.AbortTimer = nil
		}
	}

	if len(foundGame.Moves) == 2 {
		if foundGame.Board.Position().Turn() == chess.White {
			timer := time.AfterFunc(foundGame.TimeWhite, func() { c.handleGameTimeout(foundGame) })
			foundGame.ClockTimer = timer
		} else {
			timer := time.AfterFunc(foundGame.TimeBlack, func() { c.handleGameTimeout(foundGame) })
			foundGame.ClockTimer = timer
		}
	} else if len(foundGame.Moves) > 2 {
		if foundGame.Board.Position().Turn() == chess.White {
			foundGame.ClockTimer.Reset(foundGame.TimeWhite)
		} else {
			foundGame.ClockTimer.Reset(foundGame.TimeBlack)
		}
	}

	// log.Println(foundGame.Board.Position().Turn())
	var cw, cb int
	var err error
	if result != "" {
		// log.Println("game ho gya over")
		etlb := int32(foundGame.TimeBlack.Milliseconds())
		etlw := int32(foundGame.TimeWhite.Milliseconds())
		cw, cb, err = c.endGame(foundGame, result, &message, int16(len(foundGame.Moves)), &etlw, &etlb)
		if err != nil {
			log.Println("error ending game with result", err)
		}
		foundGame.ClockTimer.Stop()

		c.sendScoreUpdateEvent(foundGame)

	}
	payload, err := json.Marshal(map[string]any{"gameID": gameID, "move": moveToSend, "Result": result, "message": message, "timeBlack": foundGame.TimeBlack.Milliseconds(), "timeWhite": foundGame.TimeWhite.Milliseconds(), "changeW": cw, "changeB": cb})
	if err != nil {
		log.Println("error marshalling new game payload")
		return nil
	}
	e := socket.Event{
		Type:    "Move_Response",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.BroadcastToRoom(e, gameID)

	if result != "" {
		c.BatchInsertMoves(foundGame)
		delete(c.GameManager.Games, gameID)
	}

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

	if len(foundGame.Moves) < 2 {
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
		c.SocketManager.SendToUserClientsInARoom(e, client.Room, other)
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
		cw, cb, err := c.endGame(foundGame, "1/2-1/2", &reason, int16(len(foundGame.Moves)), &etlw, &etlb)
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
		c.sendScoreUpdateEvent(foundGame)

		c.BatchInsertMoves(foundGame)
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

	if len(foundGame.Moves) < 2 {
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
	cw, cb, err := c.endGame(foundGame, result, &reason, int16(len(foundGame.Moves)), &etlw, &etlb)
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

	c.sendScoreUpdateEvent(foundGame)

	c.BatchInsertMoves(foundGame)
	delete(c.GameManager.Games, gameID)
	return nil
}
