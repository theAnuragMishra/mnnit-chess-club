package control

type movePayload struct {
	MoveStr string `json:"MoveStr"`
	Orig    string `json:"orig"`
	Dest    string `json:"dest"`
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

type ChatPayload struct {
	Sender           int32  `json:"sender"`
	Receiver         int32  `json:"receiver"`
	SenderUsername   string `json:"senderUsername"`
	ReceiverUsername string `json:"receiverUsername"`
	Text             string `json:"text"`
	GameID           int32  `json:"gameID"`
}
