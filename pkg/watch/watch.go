// xref: k8s.io/apimachinery/pkg/watch

package watch

type EventType int

const (
	internalRunFunctionMarker EventType = iota // internal do function

	// Invalid
	EventTypeInvalid

	// Player Event
	PlayerPostSmallBlind
	PlayerPostBigBlind
	PlayerCheck
	PlayerFold
	PlayerBet
	PlayerCall
	PlayerRaise
	PlayerAllIn

	// Deal Event
	DealerShuffle
	DealerHoleCards
	DealerFlopCards
	DealerTurnCard
	DealerRiverCard
	DealerBurnCard

	// Round Event
	RoundStart
	RoundPreFlop
	RoundFlop
	RoundTurn
	RoundRiver
	RoundShowdown
	RoundEnd

	// Table Event
	TablePlayerJoin
	TablePlayerLeave
)

type Interface interface {
	Stop()
	Watch() <-chan Event
}

var (
	DefaultChanSize int32 = 100
)

type Event struct {
	Type   EventType
	Object any
}

type emptyWatch chan Event

func NewEmptyWatch() Interface {
	ch := make(chan Event)
	close(ch)
	return emptyWatch(ch)
}

// Stop implements Interface
func (w emptyWatch) Stop() {
}

// ResultChan implements Interface
func (w emptyWatch) Watch() <-chan Event {
	return chan Event(w)
}
