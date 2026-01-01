package player

type ActionType int

const (
	ActionInvalid ActionType = iota
	ActionCheck
	ActionFold
	ActionBet
	ActionCall
	ActionRaise
	ActionAllIn
)

func (at ActionType) String() string {
	switch at {
	case ActionCheck:
		return "Check"
	case ActionFold:
		return "Fold"
	case ActionBet:
		return "Bet"
	case ActionCall:
		return "Call"
	case ActionRaise:
		return "Raise"
	case ActionAllIn:
		return "AllIn"
	default:
		return "Invalid"
	}
}

func (at ActionType) IntoEventType() EventType {
	switch at {
	case ActionCheck:
		return EventCheck
	case ActionFold:
		return EventFold
	case ActionBet:
		return EventBet
	case ActionCall:
		return EventCall
	case ActionRaise:
		return EventRaise
	case ActionAllIn:
		return EventAllIn
	default:
		return EventInvalid
	}
}

func (at ActionType) ToStatus() StatusType {
	switch at {
	case ActionCheck, ActionBet, ActionRaise, ActionCall:
		return StatusWaitingToAct
	case ActionFold:
		return StatusFolded
	case ActionAllIn:
		return StatusAllIn
	default:
		return StatusReady
	}
}

type Action struct {
	Type ActionType

	// for Bet, Call, Raise
	Chips int
}
