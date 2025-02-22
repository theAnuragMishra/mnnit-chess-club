package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func (c *Controller) InitGame(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(event)
	var user userPayload
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	session := r.Context().Value(auth.MiddlewareSentSession).(database.Session)

	log.Println("inside init game :)")

	// if there's no pending game, create one, else, add the player to the game
	if c.GameManager.PendingUserID == 0 {
		log.Println("no pending game, creating...")
		c.GameManager.PendingUserName = user.Username
		c.GameManager.PendingUserID = session.UserID

		//newGame := game.NewGame(client.UserID, user.Username)
		//c.GameManager.PendingGameID = newGame.ID
		//c.GameManager.Games = append(c.GameManager.Games, newGame)
		// payload := map[string]interface{}{
		// 	"GameID":          newGame.ID,
		// 	"player1":         client.UserID,
		// 	"player2":         "",
		// 	"player1username": user.Username,
		// 	"player2username": "",
		// }
		// rawPayload, err := json.Marshal(payload)
		// if err != nil {
		// 	log.Println("error marshalling new game payload")
		// 	return nil
		// }
		//
		// e := socket.Event{
		// 	Type:    "Init_Game",
		// 	Payload: json.RawMessage(rawPayload),
		// }
		// client.Send(e)
	} else {

		if c.GameManager.PendingUserName == user.Username {
			utils.RespondWithError(w, http.StatusBadRequest, "You can't play both sides")
			return
		}
		createdGame := game.NewGame(c.GameManager.PendingUserID, session.UserID)

		id, err := c.Queries.CreateGame(context.Background(), database.CreateGameParams{
			WhiteID:       &c.GameManager.PendingUserID,
			BlackID:       &session.UserID,
			WhiteUsername: c.GameManager.PendingUserName,
			BlackUsername: user.Username,
			Fen:           createdGame.Board.FEN(),
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		createdGame.ID = id
		c.GameManager.Games = append(c.GameManager.Games, createdGame)

		thisClient := c.SocketManager.FindClientByUserID(session.UserID)
		otherClient := c.SocketManager.FindClientByUserID(c.GameManager.PendingUserID)
		if thisClient == nil || otherClient == nil {
			log.Println("player not found")
			return
		}

		payload := map[string]interface{}{
			"GameID":          createdGame.ID,
			"player1":         otherClient.UserID,
			"player2":         session.UserID,
			"player1username": c.GameManager.PendingUserName,
			"player2username": user.Username,
		}

		c.GameManager.PendingUserID = 0
		c.GameManager.PendingUserName = ""

		rawPayload, err := json.Marshal(payload)
		if err != nil {
			log.Println("error marshalling new createdGame payload")
			return
		}

		e := socket.Event{
			Type:    "Init_Game",
			Payload: json.RawMessage(rawPayload),
		}

		thisClient.Send(e)
		otherClient.Send(e)

	}

	// fmt.Println(c.GameManager)
}

func Move(c *Controller, event socket.Event, client *socket.Client) error {
	// fmt.Println(string(event.Payload))
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

	payload, err := json.Marshal(map[string]interface{}{"move": insertedMove, "Result": result, "message": message})
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

func (c *Controller) WriteGameInfo(w http.ResponseWriter, r *http.Request) {
	gameIDStr, err := strconv.ParseInt(chi.URLParam(r, "gameID"), 10, 32)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
	}
	gameID := int32(gameIDStr)
	// fmt.Println(gameID)
	foundGame, err := c.Queries.GetGameInfo(r.Context(), gameID)
	if err != nil {
		log.Println(err)

		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
	}
	moves, err := c.Queries.GetGameMoves(r.Context(), gameID)
	if err != nil {
		log.Println(err)
	}

	response := GameResponse{
		Game:  foundGame,
		Moves: moves,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
