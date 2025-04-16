package control

type movePayload struct {
	MoveStr string `json:"MoveStr"`
	Orig    string `json:"orig"`
	Dest    string `json:"dest"`
	GameID  string `json:"GameId"`
}

type InitGamePayload struct {
	TimeControl string `json:"timeControl"`
}

type AcceptChallengePayload struct {
	GameID string `json:"GameId"`
}

type UserNamePayload struct {
	Username string `json:"username"`
}

type RoomPayload struct {
	RoomID string `json:"room"`
}

type GameResponse struct {
	Game      any   `json:"game"`
	Moves     any   `json:"moves"`
	TimeBlack int32 `json:"timeBlack"`
	TimeWhite int32 `json:"timeWhite"`
}

type ChatPayload struct {
	Text string `json:"text"`
}

type DRPayload struct {
	PlayerID int32  `json:"playerID"`
	GameID   string `json:"gameID"`
}
