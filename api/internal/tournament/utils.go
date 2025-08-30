package tournament

func (t *Tournament) UpdateSnapShot() {
	m := make(map[int32]Player, len(t.players))
	for k, v := range t.players {
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
	t.PlayersSnapShot.Store(m)
}

func (t *Tournament) playerSnapshot(id int32) PlayerSnapShot {
	player := t.players[id]
	return PlayerSnapShot{
		ID:     player.ID,
		Score:  player.Score,
		Scores: player.Scores,
		Rating: player.Rating,
		Streak: player.Streak,
	}
}
