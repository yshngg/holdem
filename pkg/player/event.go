package player

import (
	"time"

	"github.com/yshngg/holdem/pkg/watch"
)

type EventAction string

const (
	EventKind string = "player"

	EventPostSmallBlind EventAction = "PostSmallBlind"
	EventPostBigBlind   EventAction = "PostBigBlind"
	EventCheck          EventAction = "Check"
	EventFold           EventAction = "Fold"
	EventBet            EventAction = "Bet"
	EventCall           EventAction = "Call"
	EventRaise          EventAction = "Raise"
	EventAllIn          EventAction = "AllIn"
)

type EventObject struct {
	ID  string
	Bet int
}

type Event struct {
	action    EventAction
	object    EventObject
	eventTime time.Time
}

func NewEvent(action EventAction, object EventObject) watch.Event {
	return Event{
		action:    action,
		object:    object,
		eventTime: time.Now(),
	}
}

func (e Event) Kind() string {
	return EventKind
}

func (e Event) Action() string {
	return string(e.action)
}

func (e Event) Related() any {
	return e.object
}

func (e Event) Time() time.Time {
	return e.eventTime
}

var _ watch.Event = Event{}
