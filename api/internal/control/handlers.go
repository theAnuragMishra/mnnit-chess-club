package control

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/auth"
	"log"
	"net/http"
	"strconv"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"

	"github.com/google/uuid"
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

	fmt.Println("inside init game :)")

	// if there's no pending game, create one, else, add the player to the game
	if c.GameManager.PendingUser == uuid.Nil {
		fmt.Println("no pending game, creating...")
		c.GameManager.PendingUser = session.UserID
		c.GameManager.PendingUserName = user.Username
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

		if c.GameManager.PendingUser == session.UserID {
			utils.RespondWithError(w, http.StatusBadRequest, "You can't play both sides")
		}

		id, err := c.Queries.CreateGame(context.Background(), database.CreateGameParams{
			WhitePlayerID: c.GameManager.PendingUser,
			BlackPlayerID: session.UserID,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		createdGame := game.NewGame(id, c.GameManager.PendingUser, c.GameManager.PendingUserName, session.UserID, user.Username)
		c.GameManager.Games = append(c.GameManager.Games, createdGame)

		thisClient := c.SocketManager.FindClientByUserID(session.UserID)
		otherClient := c.SocketManager.FindClientByUserID(c.GameManager.PendingUser)
		if thisClient == nil || otherClient == nil {
			fmt.Println("player not found")
			return
		}

		payload := map[string]interface{}{
			"GameID":          createdGame.ID,
			"player1":         otherClient.UserID,
			"player2":         session.UserID,
			"player1username": c.GameManager.PendingUserName,
			"player2username": user.Username,
		}

		c.GameManager.PendingUser = uuid.Nil
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
	fmt.Println(string(event.Payload))
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

	if foundGame.WhitePlayerID == uuid.Nil || foundGame.BlackPlayerID == uuid.Nil {
		return errors.New("game not started")
	}

	moveResult := foundGame.MakeMove(client.UserID, move.MoveStr)

	payload, err := json.Marshal(map[string]interface{}{"Result": moveResult})
	if err != nil {
		log.Println("error marshalling new game payload")
		return nil
	}
	e := socket.Event{
		Type:    "Result_Alert",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.Broadcast(e)

	return nil
}

func (c *Controller) WriteGameInfo(w http.ResponseWriter, r *http.Request) {
	gameIDStr, err := strconv.ParseInt(chi.URLParam(r, "gameID"), 10, 32)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
	}
	gameID := int32(gameIDStr)
	fmt.Println(gameID)
	foundGame, err := c.Queries.GetGameInfo(r.Context(), gameID)
	if err != nil {
		fmt.Println(err)

		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID")
	}
	utils.RespondWithJSON(w, http.StatusOK, foundGame)
}
