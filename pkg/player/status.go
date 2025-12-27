package player

type Status int

const (
	_ Status = iota
	StatusFold
	StatusCheck
	StatusBet
	StatusRaise
	StatusAllIn
	StatusReady
	StatusActive // who is taking action
)

func (s Status) String() string {
	switch s {
	case StatusFold:
		return "Fold"
	case StatusCheck:
		return "Check"
	case StatusBet:
		return "Bet"
	case StatusRaise:
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
