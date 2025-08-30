package game

import (
	"crypto/rand"
	"time"

	"github.com/notnil/chess"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateUniqueID(length int) (string, error) {
	code := make([]byte, length)
	charsetLen := byte(len(charset))
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	for i, b := range randomBytes {
		code[i] = charset[b%charsetLen]
	}

	return string(code), nil
}

func (g *Game) setUpAbort() {
	var t time.Duration
	if g.BaseTime >= 20*time.Second {
		t = time.Second * 20
	} else if g.BaseTime >= 10*time.Second {
		t = time.Second * 10
	} else {
		t = g.BaseTime
	}
	timer := time.AfterFunc(t, func() { g.inbox <- Abort{} })
	g.st.AbortTimer = timer
}

func (g *Game) snapshot() SnapShot {
	if len(g.st.Moves) >= 2 {
		timePassed := time.Since(g.st.LastMoveTime)
		if g.st.Board.Position().Turn() == chess.White {
			g.st.TimeWhite = max(g.st.TimeWhite-timePassed, 0)
		} else {
			g.st.TimeBlack = max(g.st.TimeBlack-timePassed, 0)
		}
		g.st.LastMoveTime = time.Now()
	}
	moves := make([]Move, len(g.st.Moves))
	copy(moves, g.st.Moves)
	out := SnapShot{
		Result:        g.st.Result,
		TimeWhite:     g.st.TimeWhite.Milliseconds(),
		TimeBlack:     g.st.TimeBlack.Milliseconds(),
		DrawOfferedBy: g.st.DrawOfferedBy,
		Moves:         moves,
		RematchOffer:  g.st.RematchOffer,
	}
	return out
}
