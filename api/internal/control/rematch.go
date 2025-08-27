package control

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

type rematchInfo struct {
	WhiteID   int32
	BlackID   int32
	BaseTime  time.Duration
	Increment time.Duration
	Offer     bool
}
type rematchManager struct {
	sync.RWMutex
	rematchCache map[string]*rematchInfo
}

func newRematchManager() *rematchManager {
	return &rematchManager{
		rematchCache: make(map[string]*rematchInfo),
	}
}

func (m *rematchManager) addRematch(id string, info *rematchInfo) {
	m.Lock()
	m.rematchCache[id] = info
	m.Unlock()
}
func (m *rematchManager) removeRematch(id string) {
	m.Lock()
	delete(m.rematchCache, id)
	m.Unlock()
}

func (m *rematchManager) getRematchByID(id string) (*rematchInfo, bool) {
	m.RLock()
	info, exists := m.rematchCache[id]
	m.RUnlock()
	return info, exists
}

func rematch(c *Controller, _ socket.Event, client *socket.Client) error {
	info, exists := c.rematchManager.getRematchByID(client.Room)
	if !exists {
		return nil
	}
	if !info.Offer {
		opp := info.WhiteID
		if info.WhiteID == client.UserID {
			opp = info.BlackID
		}
		e := socket.Event{Type: "rematchOffer", Payload: json.RawMessage("[]")}
		c.socketManager.SendToUserClientsInARoom(e, client.Room, opp)
		info.Offer = true
		return nil
	}
	c.rematchManager.removeRematch(client.Room)
	id, err := c.generateUniqueGameID()
	if err != nil {
		return err
	}
	rating1, err := c.queries.GetUserRating(context.Background(), info.BlackID)
	if err != nil {
		return fmt.Errorf("error getting rating for user %d: %w", info.BlackID, err)
	}
	rating2, err := c.queries.GetUserRating(context.Background(), info.WhiteID)
	if err != nil {
		return fmt.Errorf("error getting rating for user %d: %w", info.WhiteID, err)
	}
	g := game.New(id, info.BaseTime, info.Increment, info.BlackID, info.WhiteID, "", c.gameRecv)
	c.gameManager.AddGame(g)
	err = c.queries.CreateGame(context.Background(), database.CreateGameParams{
		ID:           id,
		BaseTime:     int32(info.BaseTime.Seconds()),
		Increment:    int32(info.Increment.Seconds()),
		WhiteID:      &info.BlackID,
		BlackID:      &info.WhiteID,
		RatingW:      int32(rating1),
		RatingB:      int32(rating2),
		TournamentID: nil,
	})
	if err != nil {
		return fmt.Errorf("error creating game: %w", err)
	}

	payload := map[string]any{"ID": id, "Type": "game"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}
	c.socketManager.SendToUserClientsInARoom(e, client.Room, info.BlackID)
	c.socketManager.SendToUserClientsInARoom(e, client.Room, info.WhiteID)
	return nil
}
