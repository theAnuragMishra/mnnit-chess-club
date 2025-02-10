package control

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func InitGame(c *Controller, event socket.Event, client *socket.Client) error {
	// fmt.Println(event)
	var user userPayload
	if err := json.Unmarshal(event.Payload, &user); err != nil {
		log.Println("Invalid user payload")
		return err
	}
	fmt.Println("inside init game :)")

	// if there's no pending game, create one, else, add the player to the game
	if c.GameManager.PendingGameId == "" {
		fmt.Println("no pending game, creating...")
		newGame := game.NewGame(client.UserID, user.Username)
		c.GameManager.PendingGameId = newGame.ID
		c.GameManager.Games = append(c.GameManager.Games, newGame)
		payload := map[string]interface{}{
			"GameID":          newGame.ID,
			"player1":         client.UserID,
			"player2":         "",
			"player1username": user.Username,
			"player2username": "",
		}
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			log.Println("error marshalling new game payload")
			return nil
		}

		e := socket.Event{
			Type:    "Init_Game",
			Payload: json.RawMessage(rawPayload),
		}
		client.Send(e)
	} else {
		var foundGame *game.Game
		for _, g := range c.GameManager.Games {
			if g.ID == c.GameManager.PendingGameId {
				foundGame = g
			}
		}

		if foundGame == nil {
			return errors.New("game not found")
		}
		if foundGame.Player1Id == client.UserID {
			payload := map[string]interface{}{
				"Message": "you can't play both sides",
			}
			rawPayload, err := json.Marshal(payload)
			if err != nil {
				log.Println("error marshalling new game payload")
				return nil
			}

			e := socket.Event{
				Type:    "Bad_Request_Warning",
				Payload: json.RawMessage(rawPayload),
			}

			client.Send(e)
			return errors.New("you can't play both sides")
		}
		fmt.Println("found pending game, ", c.GameManager.PendingGameId)
		foundGame.Player2Id = client.UserID
		foundGame.Player2Username = user.Username
		c.GameManager.PendingGameId = ""

		otherClient := c.SocketManager.FindClientByUserID(foundGame.Player1Id)
		if otherClient == nil {
			return errors.New("player1 not found")
		}

		payload := map[string]interface{}{
			"GameID":          foundGame.ID,
			"player1":         otherClient.UserID,
			"player2":         client.UserID,
			"player1username": foundGame.Player1Username,
			"player2username": user.Username,
		}
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			log.Println("error marshalling new game payload")
			return nil
		}

		e := socket.Event{
			Type:    "Init_Game",
			Payload: json.RawMessage(rawPayload),
		}
		client.Send(e)

		otherClient.Send(e)

	}

	// fmt.Println(c.GameManager)

	return nil
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

	if foundGame.Player1Id == uuid.Nil || foundGame.Player2Id == uuid.Nil {
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
