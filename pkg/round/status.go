package round

type Status int

const (
	StatusInvalid Status = iota
	StatusReady
	StatusStart
	StatusPreFlop
	StatusFlop
	StatusTurn
	StatusRiver
	StatusShowdown
	StatusEnd
)

func (s Status) String() string {
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

func (s Status) After(other Status) bool {
	return s > other
}

func (s Status) Before(other Status) bool {
	return s < other
}

func (s Status) Next() Status {
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

func (s Status) Previous() Status {
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
