package tournament

import (
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
)

type SnapShot struct {
	ID          string
	Name        string
	Players     map[int32]*Player
	StartTime   time.Time
	Duration    int32
	TimeControl game.TimeControl
	CreatedBy   int32
	Creator     string
}
type PlayerSnapShot struct {
	ID     int32
	Score  int32
	Scores []int16
	Rating float64
	Streak int32
}
type UpdatedPlayerSnapShots struct {
	Player1 PlayerSnapShot
	Player2 PlayerSnapShot
}
type EndPlayer struct {
	ID     int32   `json:"id"`
	Score  int32   `json:"score"`
	Scores []int16 `json:"scores"`
	Streak int32   `json:"streak"`
}
