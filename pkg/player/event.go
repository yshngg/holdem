package player

import (
	"github.com/yshngg/holdem/pkg/watch"
)

type EventType = watch.EventType

const (
	EventTypeUndefined  EventType = watch.EventTypeUndefined
	EventPostSmallBlind           = watch.PlayerPostSmallBlind
	EventPostBigBlind             = watch.PlayerPostBigBlind
	EventCheck                    = watch.PlayerCheck
	EventFold                     = watch.PlayerFold
	EventBet                      = watch.PlayerBet
	EventCall                     = watch.PlayerCall
	EventRaise                    = watch.PlayerRaise
	EventAllIn                    = watch.PlayerAllIn
)

type EventObject struct {
	ID    string
	Bet   int
	Raise int
}
