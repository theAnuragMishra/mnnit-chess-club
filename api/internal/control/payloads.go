package control

type movePayload struct {
	MoveStr string `json:"MoveStr"`
	Orig    string `json:"orig"`
	Dest    string `json:"dest"`
	GameID  int32  `json:"GameId"`
}

type timeupPayload struct {
	GameID int32 `json:"gameID"`
}

type InitGamePayload struct {
	Username  string `json:"username"`
	TimerCode int16  `json:"timerCode"`
}

type UserNamePayload struct {
	Username string `json:"username"`
}

type GameResponse struct {
	Game      any   `json:"game"`
	Moves     any   `json:"moves"`
	Ongoing   bool  `json:"ongoing"`
	TimeBlack int32 `json:"timeBlack"`
	TimeWhite int32 `json:"timeWhite"`
}

type ChatPayload struct {
	Sender           int32  `json:"sender"`
	Receiver         int32  `json:"receiver"`
	SenderUsername   string `json:"senderUsername"`
	ReceiverUsername string `json:"receiverUsername"`
	Text             string `json:"text"`
	GameID           int32  `json:"gameID"`
}

type DRPayload struct {
	PlayerID int32 `json:"playerID"`
	GameID   int32 `json:"gameID"`
}
