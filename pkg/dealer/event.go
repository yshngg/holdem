package dealer

import (
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
}
