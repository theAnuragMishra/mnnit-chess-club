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

func (t *Tournament) PlayerSnapshot(id int32) PlayerSnapShot {
	player := t.Players[id]
	return PlayerSnapShot{
		ID:     player.ID,
		Score:  player.Score,
		Scores: player.Scores,
		Rating: player.Rating,
		Streak: player.Streak,
	}
}

func (t *Tournament) SnapshotPlayers() map[int32]Player {
	m := make(map[int32]Player, len(t.Players))
	for k, v := range t.Players {
		scores := make([]int, len(v.Scores))
		copy(scores, v.Scores)
		m[k] = Player{
			ID:              v.ID,
			IsActive:        v.IsActive,
			Score:           v.Score,
			Scores:          v.Scores,
			Rating:          v.Rating,
			Streak:          v.Streak,
			Opponents:       v.Opponents,
			LastPlayedColor: v.LastPlayedColor,
		}
	}
	return m
}
