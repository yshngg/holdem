package table

import (
	"time"

	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/watch"
)

type EventType = watch.EventType

const (
	EventPlayerJoin  = watch.TablePlayerJoin
	EventPlayerLeave = watch.TablePlayerLeave
)

type EventObject struct {
	player    *player.Player
	timestamp time.Time
}

type Event struct {
	Type   EventType
	Object any
}

func newEvent(eventType EventType, p *player.Player) Event {
	return Event{
		Type: eventType,
		Object: EventObject{
			player:    p,
			timestamp: time.Now(),
		},
	}
}
