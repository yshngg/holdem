package round

import (
	"time"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/watch"
)

type EventType = watch.EventType

const (
	EventStart    EventType = watch.RoundStart
	EventPreFlop            = watch.RoundPreFlop
	EventFlop               = watch.RoundFlop
	EventTurn               = watch.RoundTurn
	EventRiver              = watch.RoundRiver
	EventShowdown           = watch.RoundShowdown
	EventEnd                = watch.RoundEnd
)

type EventObject struct {
	CommunityCards []card.Card
	Players        []player.Player
	Timestamp      time.Time
}
