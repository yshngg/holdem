package player

type StatusType int

const (
	StatusReady        StatusType = iota // ready to start a new round
	StatusWaitingToAct                   // wait for next action
	StatusTakingAction                   // be taking action (include check, bet, raise, call)
	StatusFolded                         // have folded, abandon any claim to the pot
	StatusAllIn                          // have bet all chips and special rule comes into play
	StatusSpectating                     // spectating the game (for future expansion)
)

func (s StatusType) String() string {
	switch s {
	case StatusReady:
		return "Ready"
	case StatusWaitingToAct:
		return "WaitingToAct"
	case StatusTakingAction:
		return "TakingAction"
	case StatusFolded:
		return "Folded"
	case StatusAllIn:
		return "AllIn"
	case StatusSpectating:
		return "Spectating"
	default:
		return "Invalid"
	}
}
