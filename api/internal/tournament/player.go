package tournament

import "github.com/notnil/chess"

type Player struct {
	Id              int32
	IsActive        bool
	Score           int32
	Scores          []int16
	Rating          float64
	Streak          int32
	Opponents       map[int32]int16
	LastPlayedColor chess.Color
}

func NewPlayer(id int32, rating float64, connected bool) *Player {
	return &Player{
		Id:        id,
		IsActive:  true,
		Score:     0,
		Scores:    make([]int16, 0),
		Rating:    rating,
		Streak:    0,
		Opponents: make(map[int32]int16),
	}
}
