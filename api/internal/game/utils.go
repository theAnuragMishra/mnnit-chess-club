package game

import (
	"crypto/rand"
	"time"
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
	timer := time.AfterFunc(t, func() {
		g.handleAbort()
	})
	g.AbortTimer = timer
}
