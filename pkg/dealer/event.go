package dealer

import (
	"github.com/yshngg/holdem/pkg/card"
	"github.com/yshngg/holdem/pkg/watch"
)

type EventType watch.EventType

const (
	EventShuffle   = watch.DealerShuffle
	EventHoleCards = watch.DealerHoleCards
	EventFlopCards = watch.DealerFlopCards
	EventTurnCard  = watch.DealerTurnCard
	EventRiverCard = watch.DealerRiverCard
	EventBurnCard  = watch.DealerBurnCard
)

type EventObject struct {
	HoleCards  [2]*card.Card
	FlopCards  [3]*card.Card
	TurnCard   *card.Card
	RiverCard  *card.Card
	BuriedCard []*card.Card
}
