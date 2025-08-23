package tournament

func (t *Tournament) snapshot() SnapShot {
	return SnapShot{
		Id:          t.Id,
		Name:        t.Name,
		Players:     t.players,
		StartTime:   t.StartTime,
		Duration:    t.Duration,
		TimeControl: t.TimeControl,
		CreatedBy:   t.CreatedBy,
		Creator:     t.Creator,
	}
}

func (t *Tournament) playerSnapshot(id int32) PlayerSnapShot {
	player := t.players[id]
	return PlayerSnapShot{
		Id:     player.Id,
		Score:  player.Score,
		Scores: player.Scores,
		Rating: player.Rating,
		Streak: player.Streak,
	}
}
