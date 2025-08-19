package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"

	"github.com/go-chi/chi/v5"
	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func (c *Controller) WriteGameInfo(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "gameID")

	foundGame, err := c.Queries.GetGameInfo(r.Context(), gameID)
	if err != nil {
		c.GameManager.RLock()
		challenge, exists := c.GameManager.PendingChallenges[gameID]
		c.GameManager.RUnlock()
		if exists {
			utils.RespondWithJSON(w, http.StatusOK, challenge)
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}
	c.GameManager.RLock()
	serverGame, exists := c.GameManager.Games[gameID]
	c.GameManager.RUnlock()
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
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"moves": moves, "game": foundGame, "timeWhite": timeWhite, "timeBlack": timeBlack, "live": false})

	} else {
		//server game response
		serverGame.Lock()
		timePassed := time.Since(serverGame.LastMoveTime)
		if serverGame.Board.Position().Turn() == chess.White {
			serverGame.TimeWhite = max(serverGame.TimeWhite-timePassed, 0)
		} else {
			serverGame.TimeBlack = max(serverGame.TimeBlack-timePassed, 0)
		}
		serverGame.LastMoveTime = time.Now()
		timeWhite := int32(serverGame.TimeWhite.Milliseconds())
		timeBlack := int32(serverGame.TimeBlack.Milliseconds())
		serverGame.Unlock()
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"moves": serverGame.Moves, "game": foundGame, "timeWhite": timeWhite, "timeBlack": timeBlack, "live": true})
	}
}

func InitGame(c *Controller, event socket.Event, client *socket.Client) error {
	var timeControl game.TimeControl
	if err := json.Unmarshal(event.Payload, &timeControl); err != nil {
		return err
	}
	//log.Println(timeControl)
	if timeControl.BaseTime <= 0 || timeControl.BaseTime > 10800 || timeControl.Increment < 0 || timeControl.Increment > 180 {
		return errors.New("invalid time control")
	}
	c.GameManager.Lock()
	p, exists := c.GameManager.PendingUsers[timeControl]
	if !exists {
		//log.Println("no pending game, creating...")
		c.GameManager.PendingUsers[timeControl] = game.PendingUser{
			ID:     client.UserID,
			Client: client,
		}
		c.GameManager.Unlock()
	} else {
		//log.Println("game found...")
		delete(c.GameManager.PendingUsers, timeControl)
		c.GameManager.Unlock()
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
		c.GameManager.Lock()
		_, err = c.createGame(id, p.ID, client.UserID, time.Duration(timeControl.BaseTime)*time.Second, time.Duration(timeControl.Increment)*time.Second, rating1, rating2, "")
		c.GameManager.Unlock()
		if err != nil {
			return err
		}
		payload := map[string]any{"ID": id, "Type": "game"}
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
	c.GameManager.Lock()
	c.GameManager.PendingChallenges[id] = game.Challenge{
		TimeControl:     timeControl,
		Creator:         client.UserID,
		CreatorUsername: client.Username,
	}
	c.GameManager.Unlock()
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
	var acceptChallengePayload GameIDPayload
	if err := json.Unmarshal(event.Payload, &acceptChallengePayload); err != nil {
		return err
	}
	c.GameManager.Lock()
	defer c.GameManager.Unlock()
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
	_, err := c.createGame(acceptChallengePayload.GameID, challenge.Creator, client.UserID, time.Duration(challenge.TimeControl.BaseTime)*time.Second, time.Duration(challenge.TimeControl.Increment)*time.Second, rating1, rating2, "")
	if err != nil {
		return err
	}
	payload := map[string]any{"ID": acceptChallengePayload.GameID, "Type": "game"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
	c.SocketManager.BroadcastToRoom(e, acceptChallengePayload.GameID)
	delete(c.GameManager.PendingChallenges, acceptChallengePayload.GameID)
	return nil
}

func Rematch(c *Controller, event socket.Event, client *socket.Client) error {
	c.GameManager.RLock()
	foundGame, exists := c.GameManager.Games[client.Room]
	c.GameManager.RUnlock()
	if !exists {
		return nil
	}
	foundGame.Lock()
	defer foundGame.Unlock()
	if !foundGame.RematchOffer {
		opp := foundGame.WhiteID
		if foundGame.WhiteID == client.UserID {
			opp = foundGame.BlackID
		}
		e := socket.Event{Type: "rematchOffer", Payload: json.RawMessage("[]")}
		c.SocketManager.SendToUserClientsInARoom(e, client.Room, opp)
		foundGame.RematchOffer = true
		return nil
	}
	id, err := c.generateUniqueGameID()
	if err != nil {
		return err
	}
	rating1, err1 := c.Queries.GetUserRating(context.Background(), foundGame.BlackID)
	rating2, err2 := c.Queries.GetUserRating(context.Background(), foundGame.WhiteID)
	if err1 != nil || err2 != nil {
		return errors.New("server error while fetching ratings")
	}
	_, err = c.createGame(id, foundGame.BlackID, foundGame.WhiteID, foundGame.BaseTime, foundGame.Increment, rating1, rating2, "")
	if err != nil {
		return err
	}

	payload := map[string]any{"ID": id, "Type": "game"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}
	c.SocketManager.SendToUserClientsInARoom(e, client.Room, foundGame.BlackID)
	c.SocketManager.SendToUserClientsInARoom(e, client.Room, foundGame.WhiteID)
	return nil
}

func Move(c *Controller, event socket.Event, client *socket.Client) error {
	// fmt.Println(string(event.Payload))

	var move movePayload
	if err := json.Unmarshal(event.Payload, &move); err != nil {
		return err
	}
	gameID := move.GameID
	c.GameManager.RLock()
	foundGame, exists := c.GameManager.Games[gameID]
	c.GameManager.RUnlock()

	if !exists {
		return errors.New("game not found")
	}
	foundGame.Lock()
	defer foundGame.Unlock()

	if foundGame.Result != 0 {
		return nil
	}

	if foundGame.Board.Position().Turn() == chess.White && client.UserID != foundGame.WhiteID {
		return errors.New("not your turn")
	}
	if foundGame.Board.Position().Turn() == chess.Black && client.UserID != foundGame.BlackID {
		return errors.New("not your turn")
	}

	err := foundGame.MakeMove(move.MoveStr)
	if err != nil {
		return err
	}

	var timeLeft int32
	if foundGame.Board.Position().Turn() == chess.Black {
		timeLeft = int32(foundGame.TimeWhite.Milliseconds())
	} else {
		timeLeft = int32(foundGame.TimeBlack.Milliseconds())
	}

	moveToSend := game.Move{
		MoveNotation: move.MoveStr,
		Orig:         move.Orig,
		Dest:         move.Dest,
		MoveFen:      foundGame.Board.FEN(),
		TimeLeft:     &timeLeft,
	}

	foundGame.Moves = append(foundGame.Moves, moveToSend)

	var cw, cb int
	var res int16
	result := foundGame.Board.Outcome()
	reason := foundGame.Board.Method().String()
	if result != "*" {
		if result == "1-0" {
			res = 1
		} else if result == "0-1" {
			res = 2
		} else {
			res = 3
		}
		// log.Println("game ho gya over")
		etlb := int32(foundGame.TimeBlack.Milliseconds())
		etlw := int32(foundGame.TimeWhite.Milliseconds())
		cw, cb, err = c.endGame(foundGame, res, &reason, int16(len(foundGame.Moves)), &etlw, &etlb)
		if err != nil {
			log.Println("error ending game with result", err)
		}
	} else {
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
	}
	payload, err := json.Marshal(map[string]any{"move": moveToSend, "Result": res, "reason": reason, "timeBlack": foundGame.TimeBlack.Milliseconds(), "timeWhite": foundGame.TimeWhite.Milliseconds(), "changeW": cw, "changeB": cb})
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

func Draw(c *Controller, event socket.Event, client *socket.Client) error {
	var draw GameIDPayload
	if err := json.Unmarshal(event.Payload, &draw); err != nil {
		return err
	}

	gameID := draw.GameID
	c.GameManager.RLock()
	foundGame, exists := c.GameManager.Games[gameID]
	c.GameManager.RUnlock()
	if !exists {
		return errors.New("game not found")
	}
	foundGame.Lock()
	defer foundGame.Unlock()
	if foundGame.Result != 0 {
		return nil
	}

	if len(foundGame.Moves) < 2 {
		return errors.New("cannot resign a game where one or both sides haven't played")
	}
	if foundGame.WhiteID != client.UserID && foundGame.BlackID != client.UserID {
		return errors.New("not one of the players")
	}

	if foundGame.DrawOfferedBy == 0 {
		foundGame.DrawOfferedBy = client.UserID
		e := socket.Event{
			Type:    "drawOffer",
			Payload: json.RawMessage("[]"),
		}

		var other int32
		if client.UserID == foundGame.BlackID {
			other = foundGame.WhiteID
		} else {
			other = foundGame.BlackID
		}
		c.SocketManager.SendToUserClientsInARoom(e, client.Room, other)
	} else if foundGame.DrawOfferedBy != client.UserID {
		reason := "Draw by mutual agreement"
		timeTaken := time.Since(foundGame.LastMoveTime)

		if foundGame.Board.Position().Turn() == chess.White {
			foundGame.TimeWhite -= timeTaken
		} else {
			foundGame.TimeBlack -= timeTaken
		}
		etlb := int32(foundGame.TimeBlack.Milliseconds())
		etlw := int32(foundGame.TimeWhite.Milliseconds())
		cw, cb, err := c.endGame(foundGame, 3, &reason, int16(len(foundGame.Moves)), &etlw, &etlb)
		if err != nil {
			log.Println("error ending game with result", err)
		}

		payload, err := json.Marshal(map[string]any{"Result": 3, "Reason": reason, "changeW": cw, "changeB": cb, "timeWhite": foundGame.TimeWhite.Milliseconds(), "timeBlack": foundGame.TimeBlack.Milliseconds()})
		if err != nil {
			return err
		}
		e := socket.Event{
			Type:    "game_end",
			Payload: json.RawMessage(payload),
		}
		c.SocketManager.BroadcastToRoom(e, gameID)
	}

	return nil
}

func Resign(c *Controller, event socket.Event, client *socket.Client) error {
	var resign GameIDPayload
	if err := json.Unmarshal(event.Payload, &resign); err != nil {
		return err
	}
	gameID := resign.GameID
	c.GameManager.RLock()
	foundGame, exists := c.GameManager.Games[gameID]
	c.GameManager.RUnlock()
	if !exists {
		return errors.New("game not found")
	}
	foundGame.Lock()
	defer foundGame.Unlock()
	if foundGame.Result != 0 {
		return nil
	}

	if len(foundGame.Moves) < 2 {
		return errors.New("cannot resign a game where one or both sides haven't played")
	}

	if foundGame.WhiteID != client.UserID && foundGame.BlackID != client.UserID {
		return errors.New("not one of the players")
	}

	var result int16
	var reason string
	if foundGame.WhiteID == client.UserID {
		result = 2
		reason = "White Resigned"
	} else {
		result = 1
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

	payload, err := json.Marshal(map[string]any{"Result": result, "Reason": reason, "changeW": cw, "changeB": cb, "timeWhite": foundGame.TimeWhite.Milliseconds(), "timeBlack": foundGame.TimeBlack.Milliseconds()})
	if err != nil {
		return err
	}
	e := socket.Event{
		Type:    "game_end",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.BroadcastToRoom(e, gameID)
	return nil
}
