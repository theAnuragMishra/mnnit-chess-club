package tournament

// messages coming in

type Msg interface {
	isMessage()
}
type GetPlayers struct {
	Reply chan map[int32]Player
}

func (GetPlayers) isMessage() {}

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
	Reply  chan Player
}

func (AddPlayer) isMessage() {}

type UpdatePlayerStatus struct {
	ID     int32
	Status bool
}

func (UpdatePlayerStatus) isMessage() {}

type UpdatePlayers struct {
	Result           int
	Player1          int32
	Player2          int32
	Rating1          float64
	Rating2          float64
	ExtraPointPlayer int32
}

func (UpdatePlayers) isMessage() {}

type EndTournament struct{}

func (EndTournament) isMessage() {}

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

type BroadCastUpdatedPlayers struct {
	TournamentID string
	Players      UpdatedPlayerSnapShots
}

func (BroadCastUpdatedPlayers) isControllerMessage() {}
