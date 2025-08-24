package tournament

import (
	"github.com/notnil/chess"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

func (t *Tournament) PairPlayers() {
	if len(t.waitingPlayers) < 2 {
		return
	}
	paired := make(map[int32]bool)
	msg := GetPairable{
		TournamentID: t.ID,
		Players:      t.waitingPlayers,
		Reply:        make(chan []*Player, 1),
	}
	t.ControllerChan <- msg
	availableToPair := <-msg.Reply
	if len(availableToPair) < 2 {
		return
	}
	for i := 0; i < len(availableToPair); i++ {
		playerA := availableToPair[i]
		if paired[playerA.ID] {
			continue
		}
		bestMatch := -1
		minScoreDiff := 1000000
		for j := i + 1; j < len(availableToPair); j++ {
			playerB := availableToPair[j]
			if paired[playerB.ID] {
				continue
			}
			currentDiff := 0
			currentDiff = utils.Abs(int(playerA.Rating) - int(playerB.Rating))
			currentDiff = utils.Abs(int(playerA.Score)-int(playerB.Score)) * 2
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
			reply := make(chan bool, 1)
			t.ControllerChan <- PairingRequest{
				TournamentID: t.ID,
				PlayerA:      playerA,
				PlayerB:      playerB,
				Reply:        reply,
			}
			ok := <-reply
			if !ok {
				continue
			}
			paired[playerA.ID] = true
			paired[playerB.ID] = true
		}
	}
	var newWaitingPlayers []*Player
	for _, player := range t.waitingPlayers {
		if !paired[player.ID] {
			newWaitingPlayers = append(newWaitingPlayers, player)
		}
	}
	t.waitingPlayers = newWaitingPlayers
}

func (t *Tournament) end() {
	players := make([]EndPlayer, 0, len(t.players))
	for _, player := range t.players {
		players = append(players, EndPlayer{
			ID:     player.ID,
			Score:  player.Score,
			Scores: player.Scores,
			Streak: player.Streak,
		})
	}
	t.ControllerChan <- EndRequest{
		TournamentID: t.ID,
		Players:      players,
	}
	t.Done <- struct{}{}
}

func (t *Tournament) handleUpdatePlayers(id1, id2 int32, result int16, r1, r2 float64) {
	p1 := t.players[id1]
	p2 := t.players[id2]
	p1.Opponents[id2] += 1
	p2.Opponents[id1] += 1
	p1.LastPlayedColor = chess.White
	p2.LastPlayedColor = chess.Black
	p1.Rating = r1
	p2.Rating = r2

	if result == 1 {
		p2.Streak = 0
		p2.Scores = append(p2.Scores, 0)
		if p1.Streak >= 2 {
			p1.Score += 4
			p1.Scores = append(p1.Scores, 4)
		} else {
			p1.Score += 2
			p1.Scores = append(p1.Scores, 2)
		}
		p1.Streak += 1
	} else if result == 2 {
		p1.Streak = 0
		p1.Scores = append(p1.Scores, 0)
		if p2.Streak >= 2 {
			p2.Score += 4
			p2.Scores = append(p2.Scores, 4)
		} else {
			p2.Score += 2
			p2.Scores = append(p2.Scores, 2)
		}
		p2.Streak += 1
	} else {
		if p1.Streak >= 2 {
			p1.Score += 2
			p1.Scores = append(p1.Scores, 2)
		} else {
			p1.Score += 1
			p1.Scores = append(p1.Scores, 1)
		}
		if p2.Streak >= 2 {
			p2.Score += 2
			p2.Scores = append(p2.Scores, 2)
		} else {
			p2.Score += 1
			p2.Scores = append(p2.Scores, 1)
		}
		p1.Streak = 0
		p2.Streak = 0
	}
	t.waitingPlayers = append(t.waitingPlayers, p2)
	t.waitingPlayers = append(t.waitingPlayers, p1)
}
