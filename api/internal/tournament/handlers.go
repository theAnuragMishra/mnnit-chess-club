package tournament

import "github.com/notnil/chess"

func (t *Tournament) UpdatePlayers(info UpdatePlayersInfo) {
	p1 := t.Players[info.Player1]
	p2 := t.Players[info.Player2]
	p1.Opponents[info.Player2] += 1
	p2.Opponents[info.Player1] += 1
	p1.LastPlayedColor = chess.White
	p2.LastPlayedColor = chess.Black
	p1.Rating = info.Rating1
	p2.Rating = info.Rating2

	p1Gets := 0
	p2Gets := 0

	if info.Result == 1 {
		p1Gets += 2
		if p1.Streak >= 2 {
			p1Gets += 2
		}
		p1.Streak += 1
		p2.Streak = 0
	} else if info.Result == 2 {
		p2Gets += 2
		if p2.Streak >= 2 {
			p2Gets += 2
		}
		p2.Streak += 1
		p1.Streak = 0
	} else {
		p1Gets += 1
		p2Gets += 1
		if p1.Streak >= 2 {
			p1Gets += 1
		}
		if p2.Streak >= 2 {
			p2Gets += 1
		}
		p1.Streak = 0
		p2.Streak = 0
	}

	if info.ExtraPointPlayer == p1.ID {
		p1Gets += 1
	} else if info.ExtraPointPlayer == p2.ID {
		p2Gets += 1
	}

	p1.Score += p1Gets
	p2.Score += p2Gets
	p1.Scores = append(p1.Scores, p1Gets)
	p2.Scores = append(p2.Scores, p2Gets)

	t.WaitingPlayers = append(t.WaitingPlayers, p2)
	t.WaitingPlayers = append(t.WaitingPlayers, p1)
}
