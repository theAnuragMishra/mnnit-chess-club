package tournament

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
)

type Tournament struct {
	ID              string
	Name            string
	players         map[int32]*Player
	StartTime       time.Time
	Duration        int32
	TimeControl     game.TimeControl
	CreatedBy       int32
	Creator         string
	waitingPlayers  []*Player
	Done            chan struct{}
	inbox           chan Msg
	ControllerChan  chan ControllerMsg
	BerserkAllowed  bool
	Status          int
	WG              sync.WaitGroup
	PlayersSnapShot atomic.Value
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
	t.UpdateSnapShot()
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
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if t.Status != 1 {
				continue
			}
			t.PairPlayers()
		case m := <-t.inbox:
			if t.Status != 1 {
				continue
			}
			switch msg := m.(type) {
			case CheckIfPlayerExists:
				_, ok := t.players[msg.ID]
				if msg.Reply != nil {
					msg.Reply <- ok
				}
			case TogglePlayerActiveMsg:
				p := t.players[msg.ID]
				p.IsActive = !p.IsActive
				if msg.Reply != nil {
					msg.Reply <- p.IsActive
				}
			case AddPlayer:
				p := NewPlayer(msg.ID, msg.Rating, true)
				t.players[msg.ID] = p
				t.waitingPlayers = append(t.waitingPlayers, p)
				if msg.Reply != nil {
					msg.Reply <- *p
				}
			case UpdatePlayerStatus:
				p := t.players[msg.ID]
				p.IsActive = msg.Status
			case UpdatePlayers:
				t.handleUpdatePlayers(msg)
			case EndTournament:
				t.end()
			}
		case <-t.Done:
			return
		}
	}
}
