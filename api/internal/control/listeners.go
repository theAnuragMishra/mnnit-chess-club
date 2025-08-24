package control

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"
)

func (c *Controller) gameReceiveListener() {
	for req := range c.gameRecv {
		go c.endGame(req)
	}
}

func (c *Controller) tournamentReceiveListener() {
	for m := range c.tournamentRecv {
		switch msg := m.(type) {
		case tournament.PairingRequest:
			t, exists := c.TournamentManager.GetTournament(msg.TournamentId)
			if !exists {
				return
			}
			id, err := c.generateUniqueGameID()
			if err != nil {
				log.Println(err)
				msg.Reply <- false
				return
			}
			g := game.New(id, time.Duration(t.TimeControl.BaseTime)*time.Second, time.Duration(t.TimeControl.Increment)*time.Second, msg.PlayerA.Id, msg.PlayerB.Id, t.Id, c.gameRecv)
			c.GameManager.AddGame(g)
			err = c.Queries.CreateGame(context.Background(), database.CreateGameParams{
				ID:           id,
				BaseTime:     t.TimeControl.BaseTime,
				Increment:    t.TimeControl.Increment,
				WhiteID:      &msg.PlayerA.Id,
				BlackID:      &msg.PlayerB.Id,
				RatingW:      int32(msg.PlayerA.Rating),
				RatingB:      int32(msg.PlayerB.Rating),
				TournamentID: &t.Id,
			})
			if err != nil {
				msg.Reply <- false
				log.Println(err)
				return
			}
			msg.Reply <- true
			payload := map[string]any{"ID": id, "Type": "game"}
			rawPayload, err := json.Marshal(payload)
			if err != nil {
				log.Println(err)
			}
			e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}

			c.SocketManager.SendToUserClientsInARoom(e, t.Id, msg.PlayerA.Id)
			c.SocketManager.SendToUserClientsInARoom(e, t.Id, msg.PlayerB.Id)
		case tournament.EndRequest:
			go c.endTournament(msg.TournamentId, msg.Players)
		case tournament.GetPairable:
			go func() {
				availableToPair := make([]*tournament.Player, 0, len(msg.Players))
				for _, player := range msg.Players {
					if player.IsActive && c.SocketManager.IsUserInARoom(msg.TournamentID, player.Id) {
						availableToPair = append(availableToPair, player)
					}
				}
				if msg.Reply != nil {
					msg.Reply <- availableToPair
				}
			}()
		}
	}
}
