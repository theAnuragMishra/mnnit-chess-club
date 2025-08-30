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
			go func() {
				t, exists := c.tournamentManager.getTournament(msg.TournamentID)
				if !exists {
					return
				}
				id, err := c.generateUniqueGameID()
				if err != nil {
					log.Println(err)
					msg.Reply <- false
					return
				}
				g := game.New(id, time.Duration(t.TimeControl.BaseTime)*time.Second, time.Duration(t.TimeControl.Increment)*time.Second, msg.PlayerA.ID, msg.PlayerB.ID, t.ID, c.gameRecv)
				c.gameManager.addGame(g)
				err = c.queries.CreateGame(context.Background(), database.CreateGameParams{
					ID:           id,
					BaseTime:     t.TimeControl.BaseTime,
					Increment:    t.TimeControl.Increment,
					WhiteID:      &msg.PlayerA.ID,
					BlackID:      &msg.PlayerB.ID,
					RatingW:      int32(msg.PlayerA.Rating),
					RatingB:      int32(msg.PlayerB.Rating),
					TournamentID: &t.ID,
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

				c.socketManager.SendToUserClientsInARoom(e, t.ID, msg.PlayerA.ID)
				c.socketManager.SendToUserClientsInARoom(e, t.ID, msg.PlayerB.ID)
			}()
		case tournament.EndRequest:
			go c.endTournament(msg.TournamentID, msg.Players)
		case tournament.GetPairable:
			go func() {
				availableToPair := make([]*tournament.Player, 0, len(msg.Players))
				for _, player := range msg.Players {
					if player.IsActive && c.socketManager.IsUserInARoom(msg.TournamentID, player.ID) {
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
