package tournament

import "github.com/notnil/chess"

type Player struct {
	ID              int32
	IsActive        bool
	Score           int
	Scores          []int
	Rating          float64
	Streak          int
	Opponents       map[int32]int
	LastPlayedColor chess.Color
}

func NewPlayer(id int32, rating float64, connected bool) *Player {
	return &Player{
		ID:        id,
		IsActive:  true,
		Score:     0,
		Scores:    make([]int, 0),
		Rating:    rating,
		Streak:    0,
		Opponents: make(map[int32]int),
	}
}
