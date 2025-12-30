package player

type Status int

const (
	StatusReady Status = iota
	StatusChecked
	StatusFolded
	StatusBetted
	StatusCalled
	StatusRaised
	StatusAllIn
)

func (s Status) String() string {
	switch s {
	case StatusReady:
		return "Ready"
	case StatusChecked:
		return "Check"
	case StatusFolded:
		return "Fold"
	case StatusBetted:
		return "Bet"
	case StatusCalled:
		return "Call"
	case StatusRaised:
		return "Raise"
	case StatusAllIn:
		return "AllIn"
	default:
		return "Invalid"
	}
}
