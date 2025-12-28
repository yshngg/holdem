package round

import (
	"time"

	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/player"
	"github.com/yshngg/holdem/pkg/watch"
)

type EventType int

const (
	Start    = watch.RoundStart
	PreFlop  = watch.RoundPreFlop
	Flop     = watch.RoundFlop
	Turn     = watch.RoundTurn
	River    = watch.RoundRiver
	Showdown = watch.RoundShowdown
	End      = watch.RoundEnd
)

type EventObject struct {
	communityCards []card.Card
	players        []player.Player
	timestamp      time.Time
}
