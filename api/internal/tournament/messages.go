package tournament

// messages coming in

type Msg interface {
	isMessage()
}

type GetState struct {
	Reply chan SnapShot
}

func (GetState) isMessage() {}

type CheckIfPlayerExists struct {
	ID    int32
	Reply chan bool
}

func (CheckIfPlayerExists) isMessage() {}

type TogglePlayerActiveMsg struct {
	ID    int32
	Reply chan bool
}

func (TogglePlayerActiveMsg) isMessage() {}

type AddPlayer struct {
	ID     int32
	Rating float64
	Reply  chan *Player
}

func (AddPlayer) isMessage() {}

type UpdatePlayerStatus struct {
	ID     int32
	Status bool
}

func (UpdatePlayerStatus) isMessage() {}

type UpdatePlayers struct {
	Result  int16
	Player1 int32
	Player2 int32
	Rating1 float64
	Rating2 float64
	Reply   chan UpdatedPlayerSnapShots
}

func (UpdatePlayers) isMessage() {}

// messages going out

type ControllerMsg interface {
	isControllerMessage()
}

type PairingRequest struct {
	TournamentID string
	PlayerA      *Player
	PlayerB      *Player
	Reply        chan bool
}

func (PairingRequest) isControllerMessage() {}

type EndRequest struct {
	TournamentID string
	Players      []EndPlayer
}

func (EndRequest) isControllerMessage() {}

type GetPairable struct {
	TournamentID string
	Players      []*Player
	Reply        chan []*Player
}

func (GetPairable) isControllerMessage() {}
