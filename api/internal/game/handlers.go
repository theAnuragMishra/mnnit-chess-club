package game

import (
	"time"

	"github.com/notnil/chess"
)

func (g *Game) handleMove(c int32, move MoveInfo) MoveResp {
	if g.st.Result != 0 {
		return MoveResp{}
	}

	if g.st.Board.Position().Turn() == chess.White && c != g.WhiteID {
		return MoveResp{}
	}
	if g.st.Board.Position().Turn() == chess.Black && c != g.BlackID {
		return MoveResp{}
	}

	if len(g.st.Moves) >= 2 {
		moveTime := time.Since(g.st.LastMoveTime)

		if g.st.Board.Position().Turn() == chess.White {
			g.st.TimeWhite -= moveTime
			if !g.BerserkWhite {
				g.st.TimeWhite += g.Increment
			}
		} else {
			g.st.TimeBlack -= moveTime
			if !g.BerserkBlack {
				g.st.TimeBlack += g.Increment
			}
		}
	}

	if err := g.st.Board.MoveStr(move.MoveStr); err != nil {
		// log.Println(err)
		return MoveResp{}
	}

	g.st.DrawOfferedBy = 0
	err := g.st.Board.Draw(chess.ThreefoldRepetition)
	if err != nil {
		g.st.Board.Draw(chess.FiftyMoveRule)
	}

	var timeLeft int32
	if g.st.Board.Position().Turn() == chess.Black {
		timeLeft = int32(g.st.TimeWhite.Milliseconds())
	} else {
		timeLeft = int32(g.st.TimeBlack.Milliseconds())
	}

	moveToSend := Move{
		MoveNotation: move.MoveStr,
		Orig:         move.Orig,
		Dest:         move.Dest,
		MoveFen:      g.st.Board.FEN(),
		TimeLeft:     &timeLeft,
	}

	g.st.Moves = append(g.st.Moves, moveToSend)

	var res int
	result := g.st.Board.Outcome()
	reason := g.st.Board.Method().String()
	if result != "*" {
		if result == "1-0" {
			res = 1
		} else if result == "0-1" {
			res = 2
		} else {
			res = 3
		}

		g.end(res, reason)

	} else {
		if g.st.AbortTimer != nil {
			if len(g.st.Moves) == 1 {
				if g.BaseTime >= time.Second*20 {
					g.st.AbortTimer.Reset(time.Second * 20)
				} else if g.BaseTime >= time.Second*10 {
					g.st.AbortTimer.Reset(time.Second * 10)
				} else {
					g.st.AbortTimer.Reset(g.BaseTime)
				}
			} else {
				g.st.AbortTimer.Stop()
				g.st.AbortTimer = nil
			}
		}

		if len(g.st.Moves) == 2 {
			timer := time.AfterFunc(g.st.TimeWhite, func() { g.inbox <- FireTimeOut{} })
			g.st.ClockTimer = timer
		} else if len(g.st.Moves) > 2 {
			if g.st.Board.Position().Turn() == chess.White {
				g.st.ClockTimer.Reset(g.st.TimeWhite)
			} else {
				g.st.ClockTimer.Reset(g.st.TimeBlack)
			}
		}
	}
	g.st.LastMoveTime = time.Now()
	return MoveResp{
		Move:      moveToSend,
		TimeBlack: g.st.TimeBlack.Milliseconds(),
		TimeWhite: g.st.TimeWhite.Milliseconds(),
	}
}

func (g *Game) end(result int, reason string) {
	g.st.Result = result
	if g.st.AbortTimer != nil {
		g.st.AbortTimer.Stop()
	}
	if g.st.ClockTimer != nil {
		g.st.ClockTimer.Stop()
	}
	timeLeftWhite := int32(g.st.TimeWhite.Milliseconds())
	timeLeftBlack := int32(g.st.TimeBlack.Milliseconds())

	var extraPointPlayer int32
	if g.st.Result == 1 && g.BerserkWhite {
		extraPointPlayer = g.WhiteID
	} else if g.st.Result == 2 && g.BerserkBlack {
		extraPointPlayer = g.BlackID
	}

	notification := EndNotification{
		Result:           result,
		Reason:           &reason,
		ID:               g.ID,
		TimeLeftWhite:    &timeLeftWhite,
		TimeLeftBlack:    &timeLeftBlack,
		WhiteID:          g.WhiteID,
		BlackID:          g.BlackID,
		Moves:            g.st.Moves,
		TournamentID:     g.TournamentID,
		BerserkBlack:     g.BerserkBlack,
		BerserkWhite:     g.BerserkWhite,
		ExtraPointPlayer: extraPointPlayer,
	}

	g.ControllerChannel <- notification
	// do the tournament stuffs
}

func (g *Game) handleDraw(c int32) int32 {
	if g.st.Result != 0 || len(g.st.Moves) < 2 || (g.WhiteID != c && g.BlackID != c) {
		return 0
	}
	if g.st.DrawOfferedBy == 0 || g.st.DrawOfferedBy == c {
		other := g.WhiteID
		if c == g.WhiteID {
			other = g.BlackID
		}
		if g.st.DrawOfferedBy != 0 {
			g.st.DrawOfferedBy = 0
		} else {
			g.st.DrawOfferedBy = c
		}
		return other
	} else {
		reason := "Draw by mutual agreement"
		timeTaken := time.Since(g.st.LastMoveTime)

		if g.st.Board.Position().Turn() == chess.White {
			g.st.TimeWhite -= timeTaken
		} else {
			g.st.TimeBlack -= timeTaken
		}
		g.end(3, reason)
		return 0
	}
}

func (g *Game) handleResign(c int32) {
	if g.st.Result != 0 {
		return
	}

	if len(g.st.Moves) < 2 {
		return
	}

	if g.WhiteID != c && g.BlackID != c {
		return
	}

	var result int
	var reason string
	if g.WhiteID == c {
		result = 2
		reason = "White Resigned"
	} else {
		result = 1
		reason = "Black Resigned"
	}

	timeTaken := time.Since(g.st.LastMoveTime)

	if g.st.Board.Position().Turn() == chess.White {
		g.st.TimeWhite -= timeTaken
	} else {
		g.st.TimeBlack -= timeTaken
	}

	g.end(result, reason)
}

func (g *Game) handleAbort() {
	if g.st.Result != 0 {
		return
	}
	g.end(4, "Game Aborted")
}

func (g *Game) handleTimeout() {
	if g.st.Result != 0 {
		return
	}
	var result int
	var reason string
	if g.st.Board.Position().Turn() == chess.White {
		g.st.TimeWhite = 0
		result = 2
		reason = "White Timeout"
	} else {
		g.st.TimeBlack = 0
		result = 1
		reason = "Black Timeout"
	}
	g.end(result, reason)
}

func (g *Game) handleBerserk(msg BerserkMsg) {
	if msg.WB == 0 {
		if !g.BerserkWhite && len(g.st.Moves) == 0 {
			g.st.TimeWhite /= 2
			g.BerserkWhite = true
			if msg.Reply != nil {
				msg.Reply <- true
			}
		} else {
			if msg.Reply != nil {
				msg.Reply <- false
			}
		}
	} else {
		if !g.BerserkBlack && len(g.st.Moves) <= 1 {
			g.st.TimeBlack /= 2
			g.BerserkBlack = true
			if msg.Reply != nil {
				msg.Reply <- true
			}
		} else {
			if msg.Reply != nil {
				msg.Reply <- false
			}
		}
	}
}
