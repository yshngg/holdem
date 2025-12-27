package player

import (
	"github.com/yshngg/holdem/pkg/watch"
)

type EventType watch.EventType

const (
	Check = watch.PlayerCheck
	Fold  = watch.PlayerFold
	Bet   = watch.PlayerBet
	Call  = watch.PlayerCall
	Raise = watch.PlayerRaise
	AllIn = watch.PlayerAllIn
)

type EventObject struct {
	player *Player
	bet    int
}
