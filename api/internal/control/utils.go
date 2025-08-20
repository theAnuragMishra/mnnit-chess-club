package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/tournament"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"

	"github.com/notnil/chess"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func (c *Controller) sendScoreUpdateEvent(g *game.Game) {
	c.TournamentManager.RLock()
	t, affectsTournament := c.TournamentManager.Tournaments[g.TournamentID]
	c.TournamentManager.RUnlock()
	if !affectsTournament {
		return
	}
	t.RLock()
	p1 := c.TournamentManager.Tournaments[g.TournamentID].Players[g.WhiteID]
	p2 := c.TournamentManager.Tournaments[g.TournamentID].Players[g.BlackID]
	t.RUnlock()
	p1Res := updateScorePlayer{
		ID:     p1.Id,
		Score:  p1.Score,
		Scores: p1.Scores,
		Rating: p1.Rating,
		Streak: p1.Streak,
	}
	p2Res := updateScorePlayer{
		ID:     p2.Id,
		Score:  p2.Score,
		Scores: p2.Scores,
		Rating: p2.Rating,
		Streak: p2.Streak,
	}
	payloadd, err := json.Marshal(map[string]any{
		"p1": p1Res,
		"p2": p2Res,
	})
	if err != nil {
		log.Println(err)
	}
	ee := socket.Event{
		Type:    "update_score",
		Payload: json.RawMessage(payloadd),
	}

	c.SocketManager.BroadcastToRoom(ee, g.TournamentID)
}

func (c *Controller) createGame(id string, p1, p2 int32, baseTime, increment time.Duration, r1 float64, r2 float64, tournamentID string) (*game.Game, error) {

	createdGame := game.NewGame(baseTime, increment, p1, p2, tournamentID)

	var tidParam *string
	if tournamentID != "" {
		tidParam = &tournamentID
	}

	err := c.Queries.CreateGame(context.Background(), database.CreateGameParams{
		ID:           id,
		BaseTime:     int32(baseTime.Seconds()),
		Increment:    int32(increment.Seconds()),
		WhiteID:      &p1,
		BlackID:      &p2,
		RatingW:      int32(r1),
		RatingB:      int32(r2),
		TournamentID: tidParam,
	})
	if err != nil {
		return nil, err
	}

	createdGame.ID = id
	c.GameManager.Games[id] = createdGame

	var t time.Duration
	if baseTime >= 20*time.Second {
		t = time.Second * 20
	} else if baseTime >= 10*time.Second {
		t = time.Second * 10
	} else {
		t = baseTime
	}

	timer := time.AfterFunc(t, func() { c.abortGame(createdGame) })
	createdGame.AbortTimer = timer
	return createdGame, nil
}

func (c *Controller) abortGame(g *game.Game) {
	g.Lock()
	defer g.Unlock()
	if g.Result != 0 {
		return
	}
	c.TournamentManager.RLock()
	t, affectsTournament := c.TournamentManager.Tournaments[g.TournamentID]
	c.TournamentManager.RUnlock()
	reason := "Game Aborted"
	var result int16 = 4

	if affectsTournament {
		if g.Board.Position().Turn() == chess.White {
			reason = "White Didn't Play"
			result = 2
			t.Lock()
			p := t.Players[g.WhiteID]
			p.IsActive = false
			t.Unlock()
		} else {
			reason = "Black Didn't Play"
			result = 1
			t.Lock()
			p := t.Players[g.BlackID]
			p.IsActive = false
			t.Unlock()
		}
	}

	etl := int32(g.BaseTime.Milliseconds())

	cw, cb, err := c.endGame(g, result, &reason, 0, &etl, &etl)
	if err != nil {
		log.Println("error ending game with result", err)
		return
	}
	payload, err := json.Marshal(map[string]any{"Result": result, "Reason": reason, "changeW": cw, "changeB": cb, "timeWhite": etl, "timeBlack": etl})
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{
		Type:    "game_end",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.BroadcastToRoom(e, g.ID)
}

func (c *Controller) handleGameTimeout(g *game.Game) {
	g.Lock()
	defer g.Unlock()
	if g.Result != 0 {
		return
	}
	var result int16
	var reason string
	if g.Board.Position().Turn() == chess.White {
		g.TimeWhite = 0
		result = 2
		reason = "White Timeout"
	} else {
		g.TimeBlack = 0
		result = 1
		reason = "Black Timeout"
	}

	etlb := int32(g.TimeBlack.Milliseconds())
	etlw := int32(g.TimeWhite.Milliseconds())

	cw, cb, err := c.endGame(g, result, &reason, int16(len(g.Moves)), &etlw, &etlb)
	if err != nil {
		log.Println("error ending game on timeout", err)
	}
	payload, err := json.Marshal(map[string]any{"Result": result, "Reason": reason, "changeW": cw, "changeB": cb, "timeWhite": etlw, "timeBlack": etlb})
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{
		Type:    "game_end",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.BroadcastToRoom(e, g.ID)
}

func (c *Controller) endGame(g *game.Game, result int16, reason *string, gameLength int16, etlw, etlb *int32) (int, int, error) {
	g.Result = result
	if result == 4 {
		err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
			Result:           result,
			ResultReason:     reason,
			ID:               g.ID,
			EndTimeLeftBlack: etlb,
			EndTimeLeftWhite: etlw,
		})
		go func() {
			c.sendScoreUpdateEvent(g)
			time.Sleep(time.Second * 30)
			e := socket.Event{Type: "GameDeleted", Payload: json.RawMessage("[]")}
			c.SocketManager.SendToUserClientsInARoom(e, g.ID, g.BlackID)
			c.SocketManager.SendToUserClientsInARoom(e, g.ID, g.WhiteID)
			c.GameManager.Lock()
			delete(c.GameManager.Games, g.ID)
			c.GameManager.Unlock()
		}()
		return 0, 0, err
	}
	c.TournamentManager.RLock()
	t, affectsTournament := c.TournamentManager.Tournaments[g.TournamentID]
	c.TournamentManager.RUnlock()
	var r float64
	if result == 1 {
		r = 1.0
	} else if result == 2 {
		r = 0.0
	} else {
		r = 0.5
	}

	p1info, err1 := c.Queries.GetRatingInfo(context.Background(), g.WhiteID)
	p2info, err2 := c.Queries.GetRatingInfo(context.Background(), g.BlackID)
	if err1 != nil || err2 != nil {
		return 0, 0, errors.New("error getting rating info")
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
	err1 = c.Queries.UpdateRating(context.Background(), database.UpdateRatingParams{
		Rating:     up1.Rating,
		Rd:         up1.RD,
		Volatility: up1.Volatility,
		ID:         g.WhiteID,
	})
	err2 = c.Queries.UpdateRating(context.Background(), database.UpdateRatingParams{
		Rating:     up2.Rating,
		Rd:         up2.RD,
		Volatility: up2.Volatility,
		ID:         g.BlackID,
	})
	if err1 != nil || err2 != nil {
		log.Println("error updating rating", err1, err2)
	}

	//update scores in tournament
	if affectsTournament {
		t.Lock()
		p1 := t.Players[g.WhiteID]
		p2 := t.Players[g.BlackID]
		p1.Opponents[g.BlackID] += 1
		p2.Opponents[g.WhiteID] += 1
		p1.LastPlayedColor = chess.White
		p2.LastPlayedColor = chess.Black
		p1.Rating = up1.Rating
		p2.Rating = up2.Rating

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
		t.WaitingPlayers = append(t.WaitingPlayers, p2)
		t.WaitingPlayers = append(t.WaitingPlayers, p1)
		t.Unlock()
	}

	cw := int32(up1.Rating - p1info.Rating)
	cb := int32(up2.Rating - p2info.Rating)

	err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
		Result:           result,
		ResultReason:     reason,
		ID:               g.ID,
		GameLength:       gameLength,
		ChangeW:          &cw,
		ChangeB:          &cb,
		EndTimeLeftWhite: etlw,
		EndTimeLeftBlack: etlb,
	})

	go func() {
		if g.ClockTimer != nil {
			g.ClockTimer.Stop()
		}
		c.BatchInsertMoves(g)
		c.sendScoreUpdateEvent(g)
		time.Sleep(time.Second * 30)
		e := socket.Event{Type: "GameDeleted", Payload: json.RawMessage("[]")}
		c.SocketManager.SendToUserClientsInARoom(e, g.ID, g.BlackID)
		c.SocketManager.SendToUserClientsInARoom(e, g.ID, g.WhiteID)
		c.GameManager.Lock()
		delete(c.GameManager.Games, g.ID)
		c.GameManager.Unlock()
	}()

	return int(up1.Rating - p1info.Rating), int(up2.Rating - p2info.Rating), err
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

func (c *Controller) BatchInsertMoves(g *game.Game) {
	moves := make([]database.InsertMovesParams, len(g.Moves))
	for i, move := range g.Moves {
		moves[i] = database.InsertMovesParams{
			GameID:       g.ID,
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

func (c *Controller) RunPairingCycle(t *tournament.Tournament, isInitial bool) {
	//log.Println("starting pairing cycle")
	t.RLock()
	if len(t.WaitingPlayers) < 2 {
		//log.Println("not enough players")
		//log.Println(t.WaitingPlayers)
		return
	}
	paired := make(map[int32]bool)
	availableToPair := make([]*tournament.Player, 0, len(t.WaitingPlayers))
	//copy(availableToPair, t.WaitingPlayers)
	for _, player := range t.WaitingPlayers {
		if c.SocketManager.IsUserInARoom(t.Id, player.Id) {
			availableToPair = append(availableToPair, player)
		}
	}

	for i := 0; i < len(availableToPair); i++ {
		//log.Println("inside i loop")
		playerA := availableToPair[i]
		if paired[playerA.Id] || !playerA.IsActive {
			continue
		}
		bestMatch := -1
		minScoreDiff := 1000000
		for j := i + 1; j < len(availableToPair); j++ {
			//log.Println("inside j loop")
			playerB := availableToPair[j]
			if paired[playerB.Id] || !playerB.IsActive {
				continue
			}
			currentDiff := 0
			if isInitial {
				currentDiff = utils.Abs(int(playerA.Rating) - int(playerB.Rating))
			} else {
				currentDiff = utils.Abs(int(playerA.Score) - int(playerB.Score))
			}
			currentDiff += int(playerA.Opponents[playerB.Id]) * 10

			if playerA.LastPlayedColor == playerB.LastPlayedColor {
				currentDiff += 20
			}

			if currentDiff < minScoreDiff {
				minScoreDiff = currentDiff
				bestMatch = j
			}
			//log.Println("out of j")
		}
		if bestMatch != -1 {
			playerB := availableToPair[bestMatch]

			id, err := c.generateUniqueGameID()
			if err != nil {
				log.Println("error generating game id", err)
				continue
			}
			r1 := playerA.Rating
			r2 := playerB.Rating
			c.GameManager.Lock()
			_, err = c.createGame(id, playerA.Id, playerB.Id, time.Duration(t.TimeControl.BaseTime)*time.Second, time.Duration(t.TimeControl.Increment)*time.Second, r1, r2, t.Id)
			c.GameManager.Unlock()
			if err != nil {
				log.Println("error creating game", err)
				continue
			}
			payload := map[string]any{"ID": id, "Type": "game"}
			rawPayload, err := json.Marshal(payload)
			if err != nil {
				log.Println(err)
			}
			e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}

			c.SocketManager.SendToUserClientsInARoom(e, t.Id, playerA.Id)
			c.SocketManager.SendToUserClientsInARoom(e, t.Id, playerB.Id)

			paired[playerA.Id] = true
			paired[playerB.Id] = true
		}
		//log.Println("out of i")
	}
	var newWaitingPlayers []*tournament.Player
	for _, player := range t.WaitingPlayers {
		if !paired[player.Id] {
			newWaitingPlayers = append(newWaitingPlayers, player)
		}
	}
	t.RUnlock()
	t.Lock()
	t.WaitingPlayers = newWaitingPlayers
	t.Unlock()
	//log.Println("out of pairing cycle")
}

func (c *Controller) StartPairingCycle(t *tournament.Tournament, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-t.Done:
				return
			case <-ticker.C:
				c.RunPairingCycle(t, false)
			}
		}
	}()
}

func (c *Controller) EndTournament(t *tournament.Tournament) {
	t.Lock()
	defer t.Unlock()
	close(t.Done)
	err := c.Queries.UpdateTournamentStatus(context.Background(), database.UpdateTournamentStatusParams{
		Status: 2,
		ID:     t.Id,
	})

	if err != nil {
		log.Println("error updating tournament status", err)
	}

	input := make([]tournamentPlayer, 0, len(t.Players))

	for _, v := range t.Players {
		input = append(input, tournamentPlayer{v.Id, v.Score, v.Scores, v.Streak})

	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		log.Println(err)
	} else {
		err = c.Queries.BatchUpdateTournamentPlayers(context.Background(), database.BatchUpdateTournamentPlayersParams{
			TournamentID: t.Id,
			PlayersInput: inputBytes,
		})
		if err != nil {
			log.Println(err)
		}
	}

	payload := map[string]any{"ID": t.Id, "Type": "tournament"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
	c.SocketManager.BroadcastToRoom(e, t.Id)
	c.TournamentManager.Lock()
	delete(c.TournamentManager.Tournaments, t.Id)
	c.TournamentManager.Unlock()
}
