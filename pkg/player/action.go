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

func (at ActionType) IntoStatus() Status {
	switch at {
	case ActionCheck:
		return StatusChecked
	case ActionFold:
		return StatusFolded
	case ActionBet:
		return StatusBetted
	case ActionCall:
		return StatusCalled
	case ActionRaise:
		return StatusRaised
	case ActionAllIn:
		return StatusAllIn
	default:
		return StatusInvalid
	}
}

type Action struct {
	Type ActionType

	// for Bet, Call, Raise
	Chips int
}
