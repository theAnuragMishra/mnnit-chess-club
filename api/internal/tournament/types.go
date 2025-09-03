package tournament

type PlayerSnapShot struct {
	ID     int32
	Score  int
	Scores []int
	Rating float64
	Streak int
}
type UpdatedPlayerSnapShots struct {
	Player1 PlayerSnapShot
	Player2 PlayerSnapShot
}
type EndPlayer struct {
	ID     int32 `json:"id"`
	Score  int   `json:"score"`
	Scores []int `json:"scores"`
	Streak int   `json:"streak"`
}

type UpdatePlayersInfo struct {
	Result           int
	Player1          int32
	Player2          int32
	Rating1          float64
	Rating2          float64
	ExtraPointPlayer int32
}
