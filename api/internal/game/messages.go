package game

// messages coming in

type Msg interface {
	isMessage()
}

type MoveInfo struct {
	MoveStr string `json:"MoveStr"`
	Orig    string `json:"orig"`
	Dest    string `json:"dest"`
	GameID  string `json:"GameID"`
}

type MoveMessage struct {
	Player int32
	Move   MoveInfo
	Reply  chan MoveResp
}

func (MoveMessage) isMessage() {}

type MoveResp struct {
	Move      Move
	TimeBlack int64
	TimeWhite int64
}

type GetState struct {
	Reply chan SnapShot
}

func (GetState) isMessage() {}

type Abort struct{}

func (Abort) isMessage() {}

type FireTimeOut struct{}

func (FireTimeOut) isMessage() {}

type DrawMsg struct {
	Player int32
	Reply  chan int32
}

func (DrawMsg) isMessage() {}

type ResignMsg struct {
	Player int32
}

func (ResignMsg) isMessage() {}

// messages going out
