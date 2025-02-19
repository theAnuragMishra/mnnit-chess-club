package control

type movePayload struct {
	MoveStr string `json:"MoveStr"`
	GameID  int32  `json:"GameId"`
}

type userPayload struct {
	Username string `json:"username"`
}

type GameResponse struct {
	Game  interface{} `json:"game"`
	Moves interface{} `json:"moves"`
}
