package player

type StatusType int

const (
	// start status
	StatusIdle  StatusType = iota
	StatusReady            // ready to start a new round, wait for dealer to deal hole cards

	// intermediate status
	StatusWaiting // wait for next action, after being dealt two hole cards
	// StatusTakingAction            // be thinking and taking action (include check, bet, raise, call)
	StatusFolded // have folded, abandon any claim to the pot
	StatusAllIn  // have bet all chips and special rule comes into play

	// final status
	StatusWon
	StatusLost
)

func (s StatusType) String() string {
	switch s {
	case StatusReady:
		return "Ready"
	case StatusWaiting:
		return "Waiting"
	// case StatusTakingAction:
	// return "TakingAction"
	case StatusFolded:
		return "Folded"
	case StatusAllIn:
		return "AllIn"
	case StatusWon:
		return "Won"
	case StatusLost:
		return "Lost"
	default:
		return "Idle"
	}
}
