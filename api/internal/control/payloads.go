package control

type movePayload struct {
	MoveStr string `json:"MoveStr"`
	GameID  string `json:"GameId"`
}

type userPayload struct {
	Username string `json:"username"`
}
