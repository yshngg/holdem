package round

import (
	"time"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/watch"
)

type EventAction string

const (
	EventKind string = "round"

	EventStart    EventAction = "Start"
	EventPreFlop  EventAction = "PreFlop"
	EventFlop     EventAction = "Flop"
	EventTurn     EventAction = "Turn"
	EventRiver    EventAction = "River"
	EventShowdown EventAction = "Showdown"
	EventEnd      EventAction = "End"
)

type PlayerInfo struct {
	ID     string
	Name   string
	Chips  int
	Status player.StatusType
}

type EventObject struct {
	Players        []PlayerInfo
	CommunityCards []*card.Card
}

type Event struct {
	action    EventAction
	object    EventObject
	eventTime time.Time
}

func NewEvent(action EventAction, players []*player.Player, communityCards ...*card.Card) watch.Event {
	infos := make([]PlayerInfo, 0, len(players))
	for _, p := range players {
		infos = append(infos, PlayerInfo{
			ID:     p.ID(),
			Name:   p.Name(),
			Chips:  p.Chips(),
			Status: p.Status(),
		})
	}
	return Event{
		action: action,
		object: EventObject{
			Players:        infos,
			CommunityCards: communityCards,
		},
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
