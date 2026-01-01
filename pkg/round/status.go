package round

type StatusType int

const (
	StatusInvalid StatusType = iota
	StatusReady
	StatusStart
	StatusPreFlop
	StatusFlop
	StatusTurn
	StatusRiver
	StatusShowdown
	StatusEnd
)

func (s StatusType) String() string {
	switch s {
	case StatusStart:
		return "Start"
	case StatusPreFlop:
		return "PreFlop"
	case StatusFlop:
		return "Flop"
	case StatusTurn:
		return "Turn"
	case StatusRiver:
		return "River"
	case StatusShowdown:
		return "Showdown"
	case StatusEnd:
		return "End"
	default:
		return "Invalid"
	}
}

func (s StatusType) After(other StatusType) bool {
	return s > other
}

func (s StatusType) Before(other StatusType) bool {
	return s < other
}

func (s StatusType) Next() StatusType {
	switch s {
	case StatusStart:
		return StatusPreFlop
	case StatusPreFlop:
		return StatusFlop
	case StatusFlop:
		return StatusTurn
	case StatusTurn:
		return StatusRiver
	case StatusRiver:
		return StatusShowdown
	case StatusShowdown:
		return StatusEnd
	default:
		return StatusInvalid
	}
}

func (s StatusType) Previous() StatusType {
	switch s {
	case StatusEnd:
		return StatusRiver
	case StatusRiver:
		return StatusTurn
	case StatusTurn:
		return StatusFlop
	case StatusFlop:
		return StatusPreFlop
	case StatusPreFlop:
		return StatusStart
	default:
		return StatusInvalid
	}
}
