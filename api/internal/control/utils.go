package control

import (
	"context"
	"encoding/json"
	"log"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func (c *Controller) abortGame(g *game.Game) {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()

	reason := "Game Aborted"
	etl := int32(g.BaseTime.Seconds())

	err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
		Result:           "aborted",
		EndTimeLeftWhite: &etl,
		EndTimeLeftBlack: &etl,
		ResultReason:     &reason,
		ID:               g.ID,
	})
	if err != nil {
		log.Println("error ending game with result", err)
		return
	}
	payload, err := json.Marshal(map[string]any{"gameID": g.ID, "Result": "Aborted", "Reason": reason})
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{
		Type:    "game_abort",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.Broadcast(e)
	delete(c.GameManager.Games, g.ID)
}
