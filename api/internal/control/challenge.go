package control

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

type challengeManager struct {
	sync.RWMutex
	pendingChallenges map[string]game.Challenge
}

func newChallengeManager() *challengeManager {
	return &challengeManager{
		pendingChallenges: make(map[string]game.Challenge),
	}
}

func (m *challengeManager) AddChallenge(id string, c game.Challenge) {
	m.Lock()
	m.pendingChallenges[id] = c
	m.Unlock()
}

func (m *challengeManager) RemoveChallenge(id string) {
	m.Lock()
	delete(m.pendingChallenges, id)
	m.Unlock()
}

func (m *challengeManager) GetChallengeByID(id string) (game.Challenge, bool) {
	m.RLock()
	c, exists := m.pendingChallenges[id]
	m.RUnlock()
	return c, exists
}

func createChallenge(c *Controller, event socket.Event, client *socket.Client) error {
	var timeControl game.TimeControl
	if err := json.Unmarshal(event.Payload, &timeControl); err != nil {
		return err
	}
	//log.Println(timeControl)
	if timeControl.BaseTime <= 0 || timeControl.BaseTime > 10800 || timeControl.Increment < 0 || timeControl.Increment > 180 {
		return errors.New("invalid time control")
	}
	id, err := c.generateUniqueGameID()
	if err != nil {
		return err
	}
	c.challengeManager.AddChallenge(id, game.Challenge{
		TimeControl:     timeControl,
		Creator:         client.UserID,
		CreatorUsername: client.Username,
	})
	payload := map[string]any{"ID": id, "Type": "game"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}
	client.Send(e)
	return nil
}

func acceptChallenge(c *Controller, event socket.Event, client *socket.Client) error {
	var acceptChallengePayload GameIDPayload
	if err := json.Unmarshal(event.Payload, &acceptChallengePayload); err != nil {
		return err
	}
	challenge, exists := c.challengeManager.GetChallengeByID(acceptChallengePayload.GameID)
	if !exists {
		return nil
	}
	if client.UserID == challenge.Creator {
		return nil
	}
	rating1, err := c.queries.GetUserRating(context.Background(), challenge.Creator)
	if err != nil {
		return fmt.Errorf("error getting rating for user %d: %w", challenge.Creator, err)
	}
	rating2, err := c.queries.GetUserRating(context.Background(), client.UserID)
	if err != nil {
		return fmt.Errorf("error getting rating for user %d: %w", client.UserID, err)
	}
	g := game.New(acceptChallengePayload.GameID, time.Duration(challenge.TimeControl.BaseTime)*time.Second, time.Duration(challenge.TimeControl.Increment)*time.Second, challenge.Creator, client.UserID, "", c.gameRecv)
	c.gameManager.AddGame(g)
	err = c.queries.CreateGame(context.Background(), database.CreateGameParams{
		ID:           g.ID,
		BaseTime:     challenge.TimeControl.BaseTime,
		Increment:    challenge.TimeControl.Increment,
		WhiteID:      &challenge.Creator,
		BlackID:      &client.UserID,
		RatingW:      int32(rating1),
		RatingB:      int32(rating2),
		TournamentID: nil,
	})
	if err != nil {
		return fmt.Errorf("error creating game: %w", err)
	}
	payload := map[string]any{"ID": acceptChallengePayload.GameID, "Type": "game"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
	c.socketManager.BroadcastToRoom(e, acceptChallengePayload.GameID)
	c.challengeManager.RemoveChallenge(acceptChallengePayload.GameID)
	return nil
}
