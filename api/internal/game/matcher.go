package game

import "sync"

var TimeControls = []TimeControl{
	{BaseTime: 60, Increment: 0},
	{BaseTime: 60, Increment: 1},
	{BaseTime: 120, Increment: 1},
	{BaseTime: 180, Increment: 0},
	{BaseTime: 180, Increment: 2},
	{BaseTime: 300, Increment: 0},
	{BaseTime: 300, Increment: 3},
	{BaseTime: 600, Increment: 0},
	{BaseTime: 600, Increment: 5},
	{BaseTime: 900, Increment: 10},
	{BaseTime: 1800, Increment: 0},
	{BaseTime: 1800, Increment: 20}}

type Matcher struct {
	sync.Mutex
	pendingUsers    map[int]int32
	userTimeControl map[int32]int
}

func NewMatcher() *Matcher {
	return &Matcher{
		pendingUsers:    make(map[int]int32),
		userTimeControl: make(map[int32]int),
	}
}

func (m *Matcher) RemoveUser(id int32) {
	tc := m.userTimeControl[id]
	delete(m.pendingUsers, tc)
	delete(m.userTimeControl, id)
}

func (m *Matcher) HandleRequest(user int32, tc int) (int32, bool) {
	m.Lock()
	defer m.Unlock()
	if prevTC, exists := m.userTimeControl[user]; exists {
		if prevTC == tc {
			delete(m.userTimeControl, user)
			delete(m.pendingUsers, tc)
			return 0, false
		}
		delete(m.userTimeControl, user)
		delete(m.pendingUsers, prevTC)
	}
	if opp, ok := m.pendingUsers[tc]; ok {
		delete(m.pendingUsers, tc)
		delete(m.userTimeControl, opp)
		return opp, true
	}
	m.userTimeControl[user] = tc
	m.pendingUsers[tc] = user
	return 0, false
}
