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
	ActionShowHoleCards
	ActionHideHoleCards
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
	case ActionShowHoleCards:
		return "ShowHoleCards"
	case ActionHideHoleCards:
		return "HideHoleCards"
	default:
		return "Invalid"
	}
}

func (at ActionType) ToStatus() StatusType {
	switch at {
	case ActionCheck, ActionBet, ActionRaise, ActionCall, ActionShowHoleCards, ActionHideHoleCards:
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
