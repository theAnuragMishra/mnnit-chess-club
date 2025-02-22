package control

type movePayload struct {
	MoveStr string `json:"MoveStr"`
	GameID  int32  `json:"GameId"`
}

type timeupPayload struct {
	Color  string `json:"color"`
	GameID int32  `json:"gameID"`
}

type InitGamePayload struct {
	Username  string `json:"username"`
	TimerCode int16  `json:"timerCode"`
}

type GameResponse struct {
	Game      interface{} `json:"game"`
	Moves     interface{} `json:"moves"`
	Ongoing   bool        `json:"ongoing"`
	TimeBlack int32       `json:"timeBlack"`
	TimeWhite int32       `json:"timeWhite"`
}
