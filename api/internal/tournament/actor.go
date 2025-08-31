package tournament

import (
	"sync"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
)

type Tournament struct {
	ID             string
	Name           string
	players        map[int32]*Player
	StartTime      time.Time
	Duration       int32
	TimeControl    game.TimeControl
	PairingTicker  *time.Ticker
	CreatedBy      int32
	Creator        string
	waitingPlayers []*Player
	Done           chan struct{}
	inbox          chan Msg
	ControllerChan chan ControllerMsg
	BerserkAllowed bool
	Status         int
	WG             sync.WaitGroup
}

func New(id, name string, duration int32, creator string, createdBy, baseTime, increment int32, initialPlayers map[int32]*Player, c chan ControllerMsg, berserkAllowed bool) *Tournament {
	timeControl := game.TimeControl{
		BaseTime:  baseTime,
		Increment: increment,
	}
	t := &Tournament{
		ID:             id,
		players:        initialPlayers,
		waitingPlayers: make([]*Player, 0, len(initialPlayers)),
		StartTime:      time.Now(),
		Duration:       duration,
		TimeControl:    timeControl,
		CreatedBy:      createdBy,
		Creator:        creator,
		Name:           name,
		Done:           make(chan struct{}),
		ControllerChan: c,
		inbox:          make(chan Msg, 256),
		BerserkAllowed: berserkAllowed,
		Status:         1,
	}
	for _, v := range initialPlayers {
		t.waitingPlayers = append(t.waitingPlayers, v)
	}
	time.AfterFunc(time.Duration(duration)*time.Second, func() { t.inbox <- EndTournament{} })
	go t.run()
	return t
}

func (t *Tournament) Inbox() chan Msg {
	return t.inbox
}

func (t *Tournament) run() {
	defer close(t.Done)
	ticker := time.NewTicker(time.Second * 10)
	t.PairingTicker = ticker
	for {
		select {
		case <-ticker.C:
			t.PairPlayers()
		case m := <-t.inbox:
			switch msg := m.(type) {
			case GetPlayers:
				if msg.Reply != nil {
					msg.Reply <- t.snapshotPlayers()
				}
			case CheckIfPlayerExists:
				_, ok := t.players[msg.ID]
				if msg.Reply != nil {
					msg.Reply <- ok
				}
			case TogglePlayerActiveMsg:
				p := t.players[msg.ID]
				if t.Status == 1 {
					p.IsActive = !p.IsActive
				}
				if msg.Reply != nil {
					msg.Reply <- p.IsActive
				}
			case AddPlayer:
				if t.Status == 1 {
					p := NewPlayer(msg.ID, msg.Rating, true)
					t.players[msg.ID] = p
					t.waitingPlayers = append(t.waitingPlayers, p)
					if msg.Reply != nil {
						msg.Reply <- *p
					}
				} else {
					msg.Reply <- Player{}
				}
			case UpdatePlayerStatus:
				if t.Status == 1 {
					p := t.players[msg.ID]
					p.IsActive = msg.Status
				}
			case UpdatePlayers:
				if t.Status == 1 {
					t.handleUpdatePlayers(msg)
				}
			case EndTournament:
				t.end()
			}
		case <-t.Done:
			return
		}
	}
}
