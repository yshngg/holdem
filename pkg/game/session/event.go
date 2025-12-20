package session

import "github.com/yshngg/holdem/pkg/player"

type EventType int

const (
	EventTypeJoin EventType = iota
	EventTypeLeave
)

type Event struct {
	_type  EventType
	player *player.Player
}

func newEvent(eventType EventType, player *player.Player) Event {
	return Event{
		_type:  eventType,
		player: player,
	}
}

func (e Event) Type() EventType {
	return e._type
}

func (e Event) Player() *player.Player {
	return e.player
}
