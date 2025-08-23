package tournament

import (
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/game"
)

type Tournament struct {
	Id             string
	Name           string
	players        map[int32]*Player
	StartTime      time.Time
	Duration       int32
	TimeControl    game.TimeControl
	CreatedBy      int32
	Creator        string
	waitingPlayers []*Player
	Done           chan struct{}
	inbox          chan Msg
	ControllerChan chan ControllerMsg
}

func New(id, name string, duration int32, creator string, createdBy, baseTime, increment int32, initialPlayers map[int32]*Player, c chan ControllerMsg) *Tournament {
	timeControl := game.TimeControl{
		BaseTime:  baseTime,
		Increment: increment,
	}
	t := &Tournament{
		Id:             id,
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
	}
	for _, v := range initialPlayers {
		t.waitingPlayers = append(t.waitingPlayers, v)
	}
	time.AfterFunc(time.Duration(duration)*time.Second, func() { t.end() })
	go t.run()
	return t
}

func (t *Tournament) Inbox() chan Msg {
	return t.inbox
}

func (t *Tournament) run() {
	defer close(t.Done)
	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			t.PairPlayers()
		case m := <-t.inbox:
			switch msg := m.(type) {
			case GetState:
				if msg.Reply != nil {
					msg.Reply <- t.snapshot()
				}
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
			case UpdatePlayerConnectionStatus:
				p := t.players[msg.ID]
				p.IsConnected = msg.Connected
			case AddPlayer:
				p := NewPlayer(msg.ID, msg.Rating, true)
				t.players[msg.ID] = p
				t.waitingPlayers = append(t.waitingPlayers, p)
				if msg.Reply != nil {
					msg.Reply <- p
				}
			case UpdatePlayerStatus:
				p := t.players[msg.ID]
				p.IsActive = msg.Status
			case UpdatePlayers:
				t.handleUpdatePlayers(msg.Player1, msg.Player2, msg.Result, msg.Rating1, msg.Rating2)
				if msg.Reply != nil {
					p1 := t.playerSnapshot(msg.Player1)
					p2 := t.playerSnapshot(msg.Player2)
					msg.Reply <- UpdatedPlayerSnapShots{
						Player1: p1,
						Player2: p2,
					}
				}
			}
		case <-t.Done:
			return
		}
	}
}
