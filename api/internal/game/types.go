package game

import (
	"time"

	"github.com/notnil/chess"
)

type TimeControl struct {
	BaseTime  int32 `json:"baseTime"`
	Increment int32 `json:"increment"`
}

type Challenge struct {
	TimeControl     TimeControl
	Creator         int32
	CreatorUsername string
}
type Move struct {
	MoveNotation string
	Orig         string
	Dest         string
	MoveFen      string
	TimeLeft     *int32
}

type SnapShot struct {
	Result        int
	TimeWhite     int64
	TimeBlack     int64
	DrawOfferedBy int32
	Moves         []Move
	RematchOffer  bool
}

type State struct {
	Result        int
	Board         *chess.Game
	TimeWhite     time.Duration
	TimeBlack     time.Duration
	DrawOfferedBy int32
	Moves         []Move
	RematchOffer  bool
	LastMoveTime  time.Time
	AbortTimer    *time.Timer
	ClockTimer    *time.Timer
}
type EndNotification struct {
	Result           int
	Method           int
	ID               string
	TimeLeftWhite    *int32
	TimeLeftBlack    *int32
	WhiteID          int32
	BlackID          int32
	Moves            []Move
	TournamentID     string
	BerserkWhite     bool
	BerserkBlack     bool
	ExtraPointPlayer int32
}
