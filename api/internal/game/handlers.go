package game

import (
	"time"

	"github.com/notnil/chess"
)

func (g *Game) HandleMove(c int32, move MoveInfo) MoveResp {
	if g.Result != 0 {
		return MoveResp{}
	}

	if g.Board.Position().Turn() == chess.White && c != g.WhiteID {
		return MoveResp{}
	}
	if g.Board.Position().Turn() == chess.Black && c != g.BlackID {
		return MoveResp{}
	}

	if len(g.Moves) >= 2 {
		moveTime := time.Since(g.LastMoveTime)

		if g.Board.Position().Turn() == chess.White {
			g.TimeWhite -= moveTime
			if !g.BerserkWhite {
				g.TimeWhite += g.Increment
			}
		} else {
			g.TimeBlack -= moveTime
			if !g.BerserkBlack {
				g.TimeBlack += g.Increment
			}
		}
	}

	if err := g.Board.MoveStr(move.MoveStr); err != nil {
		// log.Println(err)
		return MoveResp{}
	}

	g.DrawOfferedBy = 0
	err := g.Board.Draw(chess.ThreefoldRepetition)
	if err != nil {
		g.Board.Draw(chess.FiftyMoveRule)
	}

	var timeLeft int32
	if g.Board.Position().Turn() == chess.Black {
		timeLeft = int32(g.TimeWhite.Milliseconds())
	} else {
		timeLeft = int32(g.TimeBlack.Milliseconds())
	}

	moveToSend := Move{
		MoveNotation: move.MoveStr,
		Orig:         move.Orig,
		Dest:         move.Dest,
		MoveFen:      g.Board.FEN(),
		TimeLeft:     &timeLeft,
	}

	g.Moves = append(g.Moves, moveToSend)

	var result int
	res := g.Board.Outcome()
	method := int(g.Board.Method())
	if res != "*" {
		if res == "1-0" {
			result = 1
		} else if res == "0-1" {
			result = 2
		} else {
			result = 3
		}

		g.end(result, method)

	} else {
		if g.AbortTimer != nil {
			if len(g.Moves) == 1 {
				if g.BaseTime >= time.Second*20 {
					g.AbortTimer.Reset(time.Second * 20)
				} else if g.BaseTime >= time.Second*10 {
					g.AbortTimer.Reset(time.Second * 10)
				} else {
					g.AbortTimer.Reset(g.BaseTime)
				}
			} else {
				g.AbortTimer.Stop()
				g.AbortTimer = nil
			}
		}

		if len(g.Moves) == 2 {
			timer := time.AfterFunc(g.TimeWhite, func() {
				g.handleTimeout()
			})
			g.ClockTimer = timer
		} else if len(g.Moves) > 2 {
			if g.Board.Position().Turn() == chess.White {
				g.ClockTimer.Reset(g.TimeWhite)
			} else {
				g.ClockTimer.Reset(g.TimeBlack)
			}
		}
	}
	g.LastMoveTime = time.Now()
	return MoveResp{
		Move:      moveToSend,
		TimeBlack: g.TimeBlack.Milliseconds(),
		TimeWhite: g.TimeWhite.Milliseconds(),
	}
}

func (g *Game) end(result int, method int) {
	g.Result = result
	if g.AbortTimer != nil {
		g.AbortTimer.Stop()
	}
	if g.ClockTimer != nil {
		g.ClockTimer.Stop()
	}
	timeLeftWhite := int32(g.TimeWhite.Milliseconds())
	timeLeftBlack := int32(g.TimeBlack.Milliseconds())

	var extraPointPlayer int32
	if g.Result == 1 && g.BerserkWhite {
		extraPointPlayer = g.WhiteID
	} else if g.Result == 2 && g.BerserkBlack {
		extraPointPlayer = g.BlackID
	}

	g.EndCallback(g, EndInfo{
		Method:           method,
		TimeLeftWhite:    &timeLeftWhite,
		TimeLeftBlack:    &timeLeftBlack,
		ExtraPointPlayer: extraPointPlayer,
	})
}

func (g *Game) HandleDraw(c int32) int32 {
	if g.Result != 0 || len(g.Moves) < 2 || (g.WhiteID != c && g.BlackID != c) {
		return 0
	}
	if g.DrawOfferedBy == 0 || g.DrawOfferedBy == c {
		other := g.WhiteID
		if c == g.WhiteID {
			other = g.BlackID
		}
		if g.DrawOfferedBy != 0 {
			g.DrawOfferedBy = 0
		} else {
			g.DrawOfferedBy = c
		}
		return other
	} else {
		timeTaken := time.Since(g.LastMoveTime)

		if g.Board.Position().Turn() == chess.White {
			g.TimeWhite -= timeTaken
		} else {
			g.TimeBlack -= timeTaken
		}
		g.end(3, 10)
		return 0
	}
}

func (g *Game) HandleResign(c int32) {
	if g.Result != 0 {
		return
	}

	if len(g.Moves) < 2 {
		return
	}

	if g.WhiteID != c && g.BlackID != c {
		return
	}

	var result int
	var method int
	if g.WhiteID == c {
		result = 2
		method = 11
	} else {
		result = 1
		method = 12
	}

	timeTaken := time.Since(g.LastMoveTime)

	if g.Board.Position().Turn() == chess.White {
		g.TimeWhite -= timeTaken
	} else {
		g.TimeBlack -= timeTaken
	}

	g.end(result, method)
}

func (g *Game) handleAbort() {
	g.Lock()
	defer g.Unlock()
	if g.Result != 0 {
		return
	}
	g.end(4, 13)
}

func (g *Game) handleTimeout() {
	g.Lock()
	defer g.Unlock()
	if g.Result != 0 {
		return
	}
	var result int
	var method int
	if g.Board.Position().Turn() == chess.White {
		g.TimeWhite = 0
		result = 2
		method = 14
	} else {
		g.TimeBlack = 0
		result = 1
		method = 15
	}
	g.end(result, method)
}

func (g *Game) HandleBerserk(wb int) bool {
	if wb == 0 {
		if !g.BerserkWhite && len(g.Moves) == 0 {
			g.TimeWhite /= 2
			g.BerserkWhite = true
			return true
		} else {
			return false
		}
	} else {
		if !g.BerserkBlack && len(g.Moves) <= 1 {
			g.TimeBlack /= 2
			g.BerserkBlack = true
			return true
		} else {
			return false
		}
	}
}
