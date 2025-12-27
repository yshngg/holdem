package option

type Option int

const (
	_ Option = iota
	Check
	Fold
	Bet
	Call
	Raise
	AllIn
)

func (o Option) String() string {
	switch o {
	case Check:
		return "Check"
	case Fold:
		return "Fold"
	case Bet:
		return "Bet"
	case Call:
		return "Call"
	case Raise:
		return "Raise"
	case AllIn:
		return "AllIn"
	default:
		return "Unknown"
	}
}
