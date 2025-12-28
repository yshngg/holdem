package player

import (
	"github.com/yshngg/holdem/pkg/watch"
)

type EventType watch.EventType

const (
	EventCheck = watch.PlayerCheck
	EventFold  = watch.PlayerFold
	EventBet   = watch.PlayerBet
	EventCall  = watch.PlayerCall
	EventRaise = watch.PlayerRaise
	EventAllIn = watch.PlayerAllIn
)

type EventObject struct {
	player *Player
	bet    int
}
