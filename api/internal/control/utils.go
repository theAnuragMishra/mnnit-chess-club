package control

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"

	"github.com/notnil/chess"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/database"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func (c *Controller) createGame(id string, p1, p2 int32, timeControl game.TimeControl, r1 float64, r2 float64) (*game.Game, error) {

	createdGame := game.NewGame(time.Duration(timeControl.BaseTime)*time.Second, time.Duration(timeControl.Increment)*time.Second, p1, p2)

	err := c.Queries.CreateGame(context.Background(), database.CreateGameParams{
		ID:        id,
		BaseTime:  timeControl.BaseTime,
		Increment: timeControl.Increment,
		WhiteID:   &p1,
		BlackID:   &p2,
		RatingW:   int32(r1),
		RatingB:   int32(r2),
	})
	if err != nil {
		return nil, err
	}

	createdGame.ID = id
	c.GameManager.Games[id] = createdGame

	var t time.Duration
	if timeControl.BaseTime >= 20 {
		t = time.Second * 20
	} else if timeControl.BaseTime >= 10 {
		t = time.Second * 10
	} else {
		t = time.Duration(timeControl.BaseTime) * time.Second
	}

	timer := time.AfterFunc(t, func() { c.abortGame(createdGame) })
	createdGame.AbortTimer = timer
	return createdGame, nil
}

func (c *Controller) abortGame(g *game.Game) {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()

	reason := "Game Aborted"
	etl := int32(g.BaseTime.Milliseconds())

	cw, cb, err := c.endGame(g.ID, "aborted", &reason, g.WhiteID, g.BlackID, 0, &etl, &etl)
	if err != nil {
		log.Println("error ending game with result", err)
		return
	}
	payload, err := json.Marshal(map[string]any{"gameID": g.ID, "Result": "aborted", "Reason": reason, "changeW": cw, "changeB": cb})
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{
		Type:    "game_abort",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.BroadcastToRoom(e, g.ID)
	delete(c.GameManager.Games, g.ID)
}

func (c *Controller) handleGameTimeout(g *game.Game) {
	c.GameManager.Lock()
	defer c.GameManager.Unlock()

	var etlw, etlb int32
	var result, reason string

	if g.Board.Position().Turn() == chess.White {
		etlw = 0
		etlb = int32(g.TimeBlack.Milliseconds())
		result = "0-1"
		reason = "White Timeout"
	} else {
		etlb = 0
		etlw = int32(g.TimeWhite.Milliseconds())
		result = "1-0"
		reason = "Black Timeout"
	}

	cw, cb, err := c.endGame(g.ID, result, &reason, g.WhiteID, g.BlackID, int16(len(g.Moves)), &etlw, &etlb)
	if err != nil {
		log.Println("error ending game on timeout", err)
	}
	payload, err := json.Marshal(map[string]any{"Result": result, "Reason": reason, "gameID": g.ID, "changeW": cw, "changeB": cb})
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{
		Type:    "timeup",
		Payload: json.RawMessage(payload),
	}
	c.SocketManager.BroadcastToRoom(e, g.ID)
	c.BatchInsertMoves(g)
	delete(c.GameManager.Games, g.ID)
}

func (c *Controller) endGame(gameID string, result string, reason *string, id1, id2 int32, gameLength int16, etlw, etlb *int32) (int, int, error) {
	if result == "aborted" {
		err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
			Result:           result,
			ResultReason:     reason,
			ID:               gameID,
			EndTimeLeftBlack: etlb,
			EndTimeLeftWhite: etlw,
		})
		return 0, 0, err
	}
	var r float64
	if result == "1-0" {
		r = 1.0
	} else if result == "0-1" {
		r = 0.0
	} else {
		r = 0.5
	}
	p1info, err1 := c.Queries.GetRatingInfo(context.Background(), id1)
	p2info, err2 := c.Queries.GetRatingInfo(context.Background(), id2)
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
		ID:         id1,
	})
	err2 = c.Queries.UpdateRating(context.Background(), database.UpdateRatingParams{
		Rating:     up2.Rating,
		Rd:         up2.RD,
		Volatility: up2.Volatility,
		ID:         id2,
	})
	if err1 != nil || err2 != nil {
		log.Println("error updating rating", err1, err2)
	}

	cw := int32(up1.Rating - p1info.Rating)
	cb := int32(up2.Rating - p2info.Rating)

	err := c.Queries.EndGameWithResult(context.Background(), database.EndGameWithResultParams{
		Result:           result,
		ResultReason:     reason,
		ID:               gameID,
		GameLength:       gameLength,
		ChangeW:          &cw,
		ChangeB:          &cb,
		EndTimeLeftWhite: etlw,
		EndTimeLeftBlack: etlb,
	})
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
	if len(t.WaitingPlayers) < 2 {
		return
	}
	paired := make(map[int32]bool)
	t.Lock()
	defer t.Unlock()
	availableToPair := make([]*tournament.Player, len(t.WaitingPlayers))
	copy(availableToPair, t.WaitingPlayers)

	for i := 0; i < len(availableToPair); i++ {
		playerA := availableToPair[i]
		clientA := c.SocketManager.FindClientByUserID(playerA.Id)
		if clientA == nil || paired[playerA.Id] || !playerA.IsActive || clientA.Room != t.Id {
			continue
		}

		bestMatch := -1
		minScoreDiff := 1000000
		for j := i + 1; j < len(availableToPair); j++ {
			playerB := availableToPair[j]
			clientB := c.SocketManager.FindClientByUserID(playerB.Id)
			if clientB == nil || paired[playerB.Id] || !playerB.IsActive || clientB.Room != t.Id {
				continue
			}
			currentDiff := 0
			if isInitial {
				currentDiff = utils.Abs(int(playerA.Rating) - int(playerB.Rating))
			} else {
				currentDiff = utils.Abs(int(playerA.Score) - int(playerB.Score))
			}
			if currentDiff < minScoreDiff {
				minScoreDiff = currentDiff
				bestMatch = j
			}
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

			createdGame, err := c.createGame(id, playerA.Id, playerB.Id, t.TimeControl, r1, r2, t.Id)

			if err != nil {
				log.Println("error creating game", err)
				continue
			}
			payload := map[string]any{"ID": createdGame.ID, "Type": "game"}
			rawPayload, err := json.Marshal(payload)
			if err != nil {
				log.Println(err)
			}
			e := socket.Event{Type: "GoTo", Payload: json.RawMessage(rawPayload)}

			whiteClient := c.SocketManager.FindClientByUserID(playerA.Id)
			blackClient := c.SocketManager.FindClientByUserID(playerB.Id)

			if whiteClient != nil {
				whiteClient.Send(e)
			}
			if blackClient != nil {
				blackClient.Send(e)
			}

			paired[playerA.Id] = true
			paired[playerB.Id] = true
		}
	}
	var newWaitingPlayers []*tournament.Player
	for _, player := range t.WaitingPlayers {
		if !paired[player.Id] {
			newWaitingPlayers = append(newWaitingPlayers, player)
		}
	}
	t.WaitingPlayers = newWaitingPlayers
}

func (c *Controller) StartPairingCycle(t *tournament.Tournament, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	go func() {
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
	close(t.Done)

	ids := make([]int32, len(t.Players))
	scores := make([]int32, len(t.Players))

	for i, player := range t.Players {
		ids[i] = player.Id
		scores[i] = player.Score
	}

	err := c.Queries.BatchUpdateScores(context.Background(), database.BatchUpdateScoresParams{
		Ids:    ids,
		Scores: scores,
	})

	if err != nil {
		log.Println(err)
	}

	payload := map[string]any{"ID": t.Id, "Type": "tournament"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	e := socket.Event{Type: "Refresh", Payload: json.RawMessage(rawPayload)}
	c.SocketManager.BroadcastToRoom(e, t.Id)
	delete(c.TournamentManager.Tournaments, t.Id)
}
