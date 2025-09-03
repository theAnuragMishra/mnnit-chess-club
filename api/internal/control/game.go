package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

type gameManager struct {
	sync.RWMutex
	games map[string]*game.Game
}

func newGameManager() *gameManager {
	return &gameManager{
		games: make(map[string]*game.Game),
	}
}

func (m *gameManager) addGame(g *game.Game) {
	m.Lock()
	m.games[g.ID] = g
	m.Unlock()
}

func (m *gameManager) removeGame(id string) {
	m.Lock()
	delete(m.games, id)
	m.Unlock()
}

func (m *gameManager) getGameByID(id string) (*game.Game, bool) {
	m.RLock()
	g, exists := m.games[id]
	m.RUnlock()
	return g, exists
}

func move(c *Controller, event socket.Event, client *socket.Client) error {
	var move game.MoveInfo
	if err := json.Unmarshal(event.Payload, &move); err != nil {
		return err
	}
	gameID := move.GameID
	g, exists := c.gameManager.getGameByID(gameID)
	if !exists {
		return errors.New("game not found")
	}
	g.Lock()
	reply := g.HandleMove(client.UserID, move)
	g.Unlock()
	if (reply == game.MoveResp{}) {
		return nil
	}
	payload, err := json.Marshal(reply)
	if err != nil {
		log.Println("error marshalling new game payload")
		return nil
	}
	e := socket.Event{
		Type:    "Move_Response",
		Payload: json.RawMessage(payload),
	}
	c.socketManager.BroadcastToRoom(e, gameID)
	return nil
}

func draw(c *Controller, event socket.Event, client *socket.Client) error {
	var draw GameIDPayload
	if err := json.Unmarshal(event.Payload, &draw); err != nil {
		return err
	}

	gameID := draw.GameID
	g, exists := c.gameManager.getGameByID(gameID)
	if !exists {
		return nil
	}
	g.Lock()
	reply := g.HandleDraw(client.UserID)
	g.Unlock()
	if reply != 0 {
		e := socket.Event{
			Type:    "drawOffer",
			Payload: json.RawMessage("[]"),
		}
		c.socketManager.SendToUserClientsInARoom(e, client.Room, reply)
	}
	return nil
}

func resign(c *Controller, event socket.Event, client *socket.Client) error {
	var resign GameIDPayload
	if err := json.Unmarshal(event.Payload, &resign); err != nil {
		return err
	}
	gameID := resign.GameID
	g, exists := c.gameManager.getGameByID(gameID)
	if !exists {
		return nil
	}
	g.Lock()
	g.HandleResign(client.UserID)
	g.Unlock()
	return nil
}

func (c *Controller) endGame(g *game.Game, info game.EndInfo) {
	t, ok := c.tournamentManager.getTournament(g.TournamentID)
	if ok {
		if g.Result == 4 {
			if len(g.Moves)%2 == 0 {
				info.Method = 16
				g.Result = 2
				t.Players[g.WhiteID].IsActive = false
			} else {
				info.Method = 17
				g.Result = 1
				t.Players[g.BlackID].IsActive = false
			}
		}
	}
	var cw, cb int32
	if g.Result != 4 {
		var r float64
		if g.Result == 1 {
			r = 1.0
		} else if g.Result == 2 {
			r = 0.0
		} else {
			r = 0.5
		}

		p1info, err := c.queries.GetRatingInfo(context.Background(), g.WhiteID)
		if err != nil {
			log.Println("error getting rating info for user ", g.WhiteID, err)
			return
		}
		p2info, err := c.queries.GetRatingInfo(context.Background(), g.BlackID)
		if err != nil {
			log.Println("error getting rating info for user ", g.BlackID, err)
			return
		}
		p1 := utils.Player{
			Rating:     p1info.Rating,
			RD:         p1info.Rd,
			Volatility: p1info.Volatility,
		}
		p2 := utils.Player{
			Rating:     p2info.Rating,
			RD:         p2info.Rd,
			Volatility: p2info.Volatility,
		}
		up1, up2 := utils.UpdateMatch(p1, p2, r)
		err = c.queries.UpdateRating(context.Background(), database.UpdateRatingParams{
			Rating:     up1.Rating,
			Rd:         up1.RD,
			Volatility: up1.Volatility,
			ID:         g.WhiteID,
		})
		if err != nil {
			log.Println("error updating rating for", g.WhiteID, err)
		}
		err = c.queries.UpdateRating(context.Background(), database.UpdateRatingParams{
			Rating:     up2.Rating,
			Rd:         up2.RD,
			Volatility: up2.Volatility,
			ID:         g.BlackID,
		})
		if err != nil {
			log.Println("error updating rating for", g.BlackID, err)
		}

		if ok {
			t.UpdatePlayers(tournament.UpdatePlayersInfo{
				Result:           g.Result,
				Player1:          g.WhiteID,
				Player2:          g.BlackID,
				Rating1:          up1.Rating,
				Rating2:          up2.Rating,
				ExtraPointPlayer: info.ExtraPointPlayer,
			})
			c.sendScoreUpdateEvent(g, t)
		}
		cw = int32(up1.Rating - p1info.Rating)
		cb = int32(up2.Rating - p2info.Rating)
	}

	err := c.queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
		Result:           int32(g.Result),
		Method:           int32(info.Method),
		ID:               g.ID,
		GameLength:       int32(len(g.Moves)),
		ChangeW:          &cw,
		ChangeB:          &cb,
		EndTimeLeftWhite: info.TimeLeftWhite,
		EndTimeLeftBlack: info.TimeLeftBlack,
		BerserkBlack:     g.BerserkBlack,
		BerserkWhite:     g.BerserkWhite,
	})
	if err != nil {
		log.Println("error ending game id", g.ID, err)
		return
	}
	payload, err := json.Marshal(map[string]any{"Result": g.Result, "Method": info.Method, "changeW": cw, "changeB": cb, "timeWhite": info.TimeLeftWhite, "timeBlack": info.TimeLeftBlack})
	if err != nil {
		log.Println(err)
		return
	}
	e := socket.Event{
		Type:    "game_end",
		Payload: json.RawMessage(payload),
	}
	c.socketManager.BroadcastToRoom(e, g.ID)
	c.rematchManager.addRematch(g.ID, &rematchInfo{
		WhiteID:   g.WhiteID,
		BlackID:   g.BlackID,
		BaseTime:  g.BaseTime,
		Increment: g.Increment,
		Offer:     false,
	})
	time.AfterFunc(time.Second*30, func() {
		e = socket.Event{Type: "GameDeleted", Payload: json.RawMessage("[]")}
		c.socketManager.SendToUserClientsInARoom(e, g.ID, g.BlackID)
		c.socketManager.SendToUserClientsInARoom(e, g.ID, g.WhiteID)
		c.rematchManager.removeRematch(g.ID)
	})
	c.batchInsertMoves(g.ID, g.Moves)
	c.gameManager.removeGame(g.ID)
}

func berserk(c *Controller, _ socket.Event, client *socket.Client) error {
	g, ok := c.gameManager.getGameByID(client.Room)
	if !ok {
		return nil
	}
	t, exists := c.tournamentManager.getTournament(g.TournamentID)
	if !exists || !t.BerserkAllowed {
		return nil
	}
	var payload map[string]any
	if client.UserID == g.WhiteID {
		g.Lock()
		success := g.HandleBerserk(0)
		g.Unlock()
		if !success {
			return nil
		}
		payload = map[string]any{"wb": 0}
	} else {
		g.Lock()
		success := g.HandleBerserk(1)
		g.Unlock()
		if !success {
			return nil
		}
		payload = map[string]any{"wb": 1}
	}

	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{
		Type:    "berserk",
		Payload: json.RawMessage(rawPayload),
	}
	c.socketManager.BroadcastToRoom(e, client.Room)
	return nil
}
