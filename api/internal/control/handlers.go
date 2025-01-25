package control

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func InitGame(c *Controller, event socket.Event, client *socket.Client) error {
	// fmt.Println(event)
	if c.GameManager.PendingGameId == "" {
		fmt.Println("no pending game, creating...")
		// newGame := &game.Game{
		// 	Player1Id: client.UserId,
		// }
		newGame := game.NewGame(client.UserId)
		c.GameManager.PendingGameId = newGame.ID
		c.GameManager.Games = append(c.GameManager.Games, newGame)
		payload := map[string]interface{}{
			"GameID": newGame.ID,
		}
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			log.Println("error marshalling new game payload")
			return nil
		}

		e := socket.Event{
			Type:    "Game_Alert",
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

		fmt.Println("found pending game, ", c.GameManager.PendingGameId)
		foundGame.Player2Id = client.UserId
		c.GameManager.PendingGameId = ""
		payload := map[string]interface{}{
			"GameID": foundGame.ID,
		}
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			log.Println("error marshalling new game payload")
			return nil
		}

		e := socket.Event{
			Type:    "Game_Alert",
			Payload: json.RawMessage(rawPayload),
		}
		client.Send(e)
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

	foundGame.MakeMove(client.UserId, move.MoveStr)

	return nil
}
