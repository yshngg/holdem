package player

type Status int

const (
	StatusInvalid Status = iota
	StatusFolded
	StatusChecked
	StatusBetted
	StatusRaised
	StatusCalled
	StatusAllIn
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
	case StatusCalled:
		return "Call"
	case StatusAllIn:
		return "AllIn"
	default:
		return "Invalid"
	}
}
