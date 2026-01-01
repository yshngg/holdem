package player

type StatusType int

const (
	StatusReady        StatusType = iota // ready to start a new round
	StatusWaiting                        // wait for next action
	StatusTakingAction                   // be taking action (include check, bet, raise, call)
	StatusFolded                         // have folded, abandon any claim to the pot
	StatusAllIn                          // have bet all chips and special rule comes into play
)

func (s StatusType) String() string {
	switch s {
	case StatusReady:
		return "Ready"
	case StatusWaiting:
		return "Waiting"
	case StatusTakingAction:
		return "TakingAction"
	case StatusFolded:
		return "Folded"
	case StatusAllIn:
		return "AllIn"
	default:
		return "Invalid"
	}
}
