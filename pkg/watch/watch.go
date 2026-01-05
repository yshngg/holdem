// xref: k8s.io/apimachinery/pkg/watch

package watch

type eventType int

func (e eventType) String() string {
	switch e {
	case internalRunFunctionMarker:
		return "internal"
	default:
		return ""
	}
}

var _ EventType = eventType(internalRunFunctionMarker)

const (
	internalRunFunctionMarker eventType = iota // internal do function

	// Undefined
	EventTypeUndefined

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

// func (e EventType) String() string {
// 	switch e {
// 	// player
// 	case PlayerPostSmallBlind:
// 		return "player post small blind"
// 	case PlayerPostBigBlind:
// 		return "player post big blind"
// 	case PlayerCheck:
// 		return "player check"
// 	case PlayerFold:
// 		return "player fold"
// 	case PlayerBet:
// 		return "player bet"
// 	case PlayerCall:
// 		return "player call"
// 	case PlayerRaise:
// 		return "player raise"
// 	case PlayerAllIn:
// 		return "player all in"

// 	// dealer
// 	case DealerShuffle:
// 		return "dealer shuffle deck"
// 	case DealerHoleCards:
// 		return "dealer deal hole cards"
// 	case DealerFlopCards:
// 		return "dealer deal flop cards"
// 	case DealerTurnCard:
// 		return "dealer deal turn card"
// 	case DealerRiverCard:
// 		return "dealer deal river card"
// 	case DealerBurnCard:
// 		return "dealer deal burn card"

// 	// round
// 	case RoundStart:
// 		return "round start"
// 	case RoundPreFlop:
// 		return "round pre-flop"
// 	case RoundFlop:
// 		return "round flop"
// 	case RoundTurn:
// 		return "round turn"
// 	case RoundRiver:
// 		return "round river"
// 	case RoundShowdown:
// 		return "showdown"
// 	case RoundEnd:
// 		return "round end"

// 	// table
// 	case TablePlayerLeave:
// 		return "table player leave"
// 	case TablePlayerJoin:
// 		return "table player join"
// 	default:
// 		return "undefined event type"
// 	}
// }

type Interface interface {
	Stop()
	Watch() <-chan Event
}

var (
	DefaultChanSize int32 = 100
)

// type Event struct {
// 	Type   EventType
// 	Object any
// }

type EventType interface {
	String() string
	Code() int
}

type Event interface {
	Type() EventType
	String() string
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
