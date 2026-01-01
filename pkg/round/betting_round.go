package round

type BettingRound int

const (
	BettingRoundInvalid BettingRound = iota
	BettingRoundPreFlop
	BettingRoundFlop
	BettingRoundTurn
	BettingRoundRiver
)

func (br BettingRound) String() string {
	switch br {
	case BettingRoundPreFlop:
		return "PreFlop"
	case BettingRoundFlop:
		return "Flop"
	case BettingRoundTurn:
		return "Turn"
	case BettingRoundRiver:
		return "River"
	default:
		return "Invalid"
	}
}
