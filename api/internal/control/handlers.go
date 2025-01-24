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
		newGame := &game.Game{
			Player1Id: client.UserId,
		}
		c.GameManager.PendingGameId = newGame.Id
		c.GameManager.Games = append(c.GameManager.Games, newGame)
	} else {
		var foundGame *game.Game
		for _, g := range c.GameManager.Games {
			if g.Id == c.GameManager.PendingGameId {
				foundGame = g
			}
		}

		if foundGame == nil {
			return errors.New("game not found")
		}
		foundGame.Player2Id = client.UserId
	}

	fmt.Println(c.GameManager)

	return nil
}

func Move(c *Controller, event socket.Event, client *socket.Client) error {
	var move movePayload
	if json.Unmarshal([]byte(event.Payload), &move) != nil {
		log.Println("Invalid move payload")
		return errors.New("invalid payload")
	}
	gameID := move.GameID

	var foundGame *game.Game
	for _, g := range c.GameManager.Games {
		if g.Id == gameID {
			foundGame = g
		}
	}

	if foundGame == nil {
		return errors.New("game not found")
	}

	foundGame.MakeMove(client.UserId, move.MoveStr)

	return nil
}
