package control

import (
	"context"
	"encoding/json"
	"log"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func (c *Controller) sendScoreUpdateEvent(p tournament.UpdatedPlayerSnapShots, room string) {
	payload, err := json.Marshal(map[string]any{
		"p1": p.Player1,
		"p2": p.Player2,
	})
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{
		Type:    "update_score",
		Payload: json.RawMessage(payload),
	}

	c.SocketManager.BroadcastToRoom(e, room)
}

func (c *Controller) batchInsertMoves(id string, m []game.Move) {
	moves := make([]database.InsertMovesParams, len(m))
	for i, move := range m {
		moves[i] = database.InsertMovesParams{
			GameID:       id,
			MoveNumber:   int32(i + 1),
			MoveNotation: move.MoveNotation,
			Orig:         move.Orig,
			Dest:         move.Dest,
			MoveFen:      move.MoveFen,
			TimeLeft:     move.TimeLeft,
		}
	}
	_, err := c.Queries.InsertMoves(context.Background(), moves)
	if err != nil {
		log.Println("error inserting moves", err)
	}
}

func (c *Controller) generateUniqueGameID() (string, error) {
	var id string
	var err error

	for {
		id, err = game.GenerateUniqueID(12)
		if err != nil {
			log.Println("error generating id:", err)
			return "", err
		}
		_, err1 := c.Queries.GetGameByID(context.Background(), id)
		_, err2 := c.Queries.GetTournamentByID(context.Background(), id)

		if err1 == nil || err2 == nil {
			log.Println("game or tournament found with id", err)
			continue
		}
		return id, nil
	}
}

func (c *Controller) generateUniqueTournamentID() (string, error) {
	var id string
	var err error

	for {
		id, err = game.GenerateUniqueID(12)
		if err != nil {
			log.Println("error generating id:", err)
			return "", err
		}
		_, err1 := c.Queries.GetTournamentByID(context.Background(), id)
		_, err2 := c.Queries.GetGameByID(context.Background(), id)

		if err1 == nil || err2 == nil {
			log.Println("game or tournament found with id", err)
			continue
		}
		return id, nil
	}
}
