package control

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

type tournamentManager struct {
	sync.RWMutex
	tournaments map[string]*tournament.Tournament
}

func newTournamentManager() *tournamentManager {
	return &tournamentManager{
		tournaments: make(map[string]*tournament.Tournament),
	}
}

func (m *tournamentManager) addTournament(t *tournament.Tournament) {
	m.Lock()
	m.tournaments[t.ID] = t
	m.Unlock()
}

func (m *tournamentManager) removeTournament(id string) {
	m.Lock()
	delete(m.tournaments, id)
	m.Unlock()
}

func (m *tournamentManager) getTournament(id string) (*tournament.Tournament, bool) {
	m.RLock()
	t, exists := m.tournaments[id]
	m.RUnlock()
	return t, exists
}

func (c *Controller) endTournament(t *tournament.Tournament) {
	t.Lock()
	defer t.Unlock()
	close(t.Done)
	err := c.queries.UpdateTournamentStatus(context.Background(), database.UpdateTournamentStatusParams{
		Status: 2,
		ID:     t.ID,
	})
	if err != nil {
		log.Println("error updating tournament status", err)
	}

	input := make([]tournament.EndPlayer, 0, len(t.Players))
	for _, player := range t.Players {
		input = append(input, tournament.EndPlayer{
			ID:     player.ID,
			Score:  player.Score,
			Scores: player.Scores,
			Streak: player.Streak,
		})
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		log.Println(err)
	} else {
		err = c.queries.BatchUpdateTournamentPlayers(context.Background(), database.BatchUpdateTournamentPlayersParams{
			TournamentID: t.ID,
			PlayersInput: inputBytes,
		})
		if err != nil {
			log.Println("error batch updating tournament players", err)
		}
	}

	payload := map[string]any{"ID": t.ID, "Type": "tournament"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
	c.socketManager.BroadcastToRoom(e, t.ID)
	c.tournamentManager.removeTournament(t.ID)
}

func (c *Controller) runTournamentPairing(t *tournament.Tournament) {
	t.RLock()
	if len(t.WaitingPlayers) < 2 {
		t.RUnlock()
		return
	}
	paired := make(map[int32]bool)

	availableToPair := make([]*tournament.Player, 0, len(t.WaitingPlayers))
	for _, player := range t.WaitingPlayers {
		if c.socketManager.IsUserInARoom(t.ID, player.ID) {
			availableToPair = append(availableToPair, player)
		}
	}
	if len(availableToPair) < 2 {
		t.RUnlock()
		return
	}
	for i := 0; i < len(availableToPair); i++ {
		playerA := availableToPair[i]
		if paired[playerA.ID] || !playerA.IsActive {
			continue
		}
		bestMatch := -1
		minScoreDiff := 1000000
		for j := i + 1; j < len(availableToPair); j++ {
			playerB := availableToPair[j]
			if paired[playerB.ID] || !playerB.IsActive {
				continue
			}
			currentDiff := utils.Abs(int(playerA.Rating) - int(playerB.Rating))
			currentDiff += utils.Abs(playerA.Score-playerB.Score) * 2
			currentDiff += int(playerA.Opponents[playerB.ID]) * 10

			if playerA.LastPlayedColor == playerB.LastPlayedColor {
				currentDiff += 20
			}

			if currentDiff < minScoreDiff {
				minScoreDiff = currentDiff
				bestMatch = j
			}
		}
		if bestMatch != -1 {
			playerB := availableToPair[bestMatch]
			//todo optimise this
			id, err := c.generateUniqueGameID()
			if err != nil {
				log.Println(err)
				continue
			}
			g := game.New(id, time.Duration(t.TimeControl.BaseTime)*time.Second, time.Duration(t.TimeControl.Increment)*time.Second, playerA.ID, playerB.ID, t.ID, c.endGame)
			c.gameManager.addGame(g)
			err = c.queries.CreateGame(context.Background(), database.CreateGameParams{
				ID:           id,
				BaseTime:     t.TimeControl.BaseTime,
				Increment:    t.TimeControl.Increment,
				WhiteID:      &playerA.ID,
				BlackID:      &playerB.ID,
				RatingW:      int32(playerA.Rating),
				RatingB:      int32(playerB.Rating),
				TournamentID: &t.ID,
			})
			if err != nil {
				log.Println(err)
				continue
			}
			payload := map[string]any{"ID": id, "Type": "game"}
			rawPayload, err := json.Marshal(payload)
			if err != nil {
				log.Println(err)
			}
			e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}

			c.socketManager.SendToUserClientsInARoom(e, t.ID, playerA.ID)
			c.socketManager.SendToUserClientsInARoom(e, t.ID, playerB.ID)
			paired[playerA.ID] = true
			paired[playerB.ID] = true
		}
	}
	var newWaitingPlayers []*tournament.Player
	for _, player := range t.WaitingPlayers {
		if !paired[player.ID] {
			newWaitingPlayers = append(newWaitingPlayers, player)
		}
	}
	t.RUnlock()
	t.Lock()
	t.WaitingPlayers = newWaitingPlayers
	t.Unlock()
}

func (c *Controller) startTournamentPairing(t *tournament.Tournament) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.runTournamentPairing(t)
				fmt.Println("tick happened")
			case <-t.Done:
				return
			}
		}
	}()
}
