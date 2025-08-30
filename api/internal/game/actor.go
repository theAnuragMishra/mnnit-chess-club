package game

import (
	"sync"
	"time"

	"github.com/notnil/chess"
)

type Game struct {
	ID                string
	BaseTime          time.Duration
	Increment         time.Duration
	inbox             chan Msg
	done              chan struct{}
	st                State
	WhiteID           int32
	BlackID           int32
	TournamentID      string
	ControllerChannel chan EndNotification
	BerserkWhite      bool
	BerserkBlack      bool
	WG                sync.WaitGroup
}

func New(id string, baseTime time.Duration, increment time.Duration, player1 int32, player2 int32, tournamentID string, c chan EndNotification) *Game {
	g := &Game{
		ID:        id,
		BaseTime:  baseTime,
		Increment: increment,
		inbox:     make(chan Msg, 256),
		done:      make(chan struct{}),
		st: State{
			Board:        chess.NewGame(),
			TimeBlack:    baseTime,
			TimeWhite:    baseTime,
			LastMoveTime: time.Now(),
		},
		WhiteID:           player1,
		BlackID:           player2,
		TournamentID:      tournamentID,
		ControllerChannel: c,
	}
	go g.run()
	return g
}

func (g *Game) Inbox() chan<- Msg {
	return g.inbox
}
func (g *Game) Done() chan<- struct{} { return g.done }
func (g *Game) run() {
	defer close(g.done)

	// wire up abort handler
	g.setUpAbort()

	for {
		select {
		case m := <-g.inbox:
			switch msg := m.(type) {
			case MoveMessage:
				resp := g.handleMove(msg.Player, msg.Move)
				if msg.Reply != nil {
					msg.Reply <- resp
				}
			case GetState:
				if msg.Reply != nil {
					msg.Reply <- g.snapshot()
				}
			case Abort:
				g.handleAbort()
			case FireTimeOut:
				g.handleTimeout()
			case DrawMsg:
				resp := g.handleDraw(msg.Player)
				if msg.Reply != nil {
					msg.Reply <- resp
				}
			case ResignMsg:
				g.handleResign(msg.Player)
			case BerserkMsg:
				g.handleBerserk(msg)
			}
		case <-g.done:
			return
		}
	}
}
