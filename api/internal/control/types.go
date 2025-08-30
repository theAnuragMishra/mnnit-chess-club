package control

import (
	"time"
)

type GameIDPayload struct {
	GameID string `json:"GameID"`
}

type UserNamePayload struct {
	Username string `json:"username"`
}

type RoomPayload struct {
	RoomID string `json:"room"`
}

type ChatPayload struct {
	Text string `json:"text"`
}
type TournamentPayload struct {
	Name           string    `json:"name"`
	BaseTime       int32     `json:"baseTime"`
	Increment      int32     `json:"increment"`
	Duration       int32     `json:"duration"`
	StartTime      time.Time `json:"startTime"`
	BerserkAllowed bool      `json:"berserkAllowed"`
}

type TournamentIDPayload struct {
	TournamentID string `json:"tournamentID"`
}
