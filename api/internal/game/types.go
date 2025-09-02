package game

type TimeControl struct {
	BaseTime  int32 `json:"baseTime"`
	Increment int32 `json:"increment"`
}
type Move struct {
	MoveNotation string
	Orig         string
	Dest         string
	MoveFen      string
	TimeLeft     *int32
}
type EndInfo struct {
	Method           int
	TimeLeftWhite    *int32
	TimeLeftBlack    *int32
	ExtraPointPlayer int32
}

type MoveInfo struct {
	MoveStr string `json:"MoveStr"`
	Orig    string `json:"orig"`
	Dest    string `json:"dest"`
	GameID  string `json:"GameID"`
}

type MoveResp struct {
	Move      Move
	TimeBlack int64
	TimeWhite int64
}
