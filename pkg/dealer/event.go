package dealer

import (
	"fmt"
	"time"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/watch"
)

type EventAction = string

const (
	EventKind string = "dealer"

	EventShuffle       EventAction = "Shuffle"
	EventDealHoleCards EventAction = "DealHoleCards"
	EventDealFlopCards EventAction = "DealFlopCards"
	EventDealTurnCard  EventAction = "DealTurnCard"
	EventDealRiverCard EventAction = "DealRiverCard"
	EventBurnCard      EventAction = "BurnCard"
)

func ToPlayer(p *player.Player) string {
	return fmt.Sprintf("player %s (id: %s)", p.Name(), p.ID())
}

func ToCommunity() string {
	return "community"
}

func ToAll() string {
	return "all"
}

type EventObject struct {
	Cards []*card.Card
	To    string
}

type Event struct {
	action    EventAction
	object    EventObject
	eventTime time.Time
}

// NewEvent instance a new dealer Event
// to must not be empty, and use ToPlayer and ToCommunity helper functions
func NewEvent(action EventAction, to string, cards ...*card.Card) watch.Event {
	e := Event{
		action: action,
		object: EventObject{
			To:    to,
			Cards: cards,
		},
		eventTime: time.Now(),
	}
	return e
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
