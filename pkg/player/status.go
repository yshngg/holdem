package player

type Status int

const (
	StatusUnknown Status = iota
	StatusFolded
	StatusChecked
	StatusBetted
	StatusRaised
	StatusCalled
	StatusAllIn
	StatusReady
	StatusActive // who is taking action
)

func (s Status) String() string {
	switch s {
	case StatusFolded:
		return "Fold"
	case StatusChecked:
		return "Check"
	case StatusBetted:
		return "Bet"
	case StatusRaised:
		return "Raise"
	case StatusAllIn:
		return "AllIn"
	case StatusReady:
		return "Ready"
	case StatusActive:
		return "Active"
	default:
		return "Unknown"
	}
}
