package watch

import "time"

type Event interface {
	Kind() string
	// Type() string
	Action() string
	Related() any
	Time() time.Time
}

type internalEventActionType string

const (
	runFunction internalEventActionType = "run function" // internal do function
)

type internalEvent struct {
	action internalEventActionType
	time   time.Time
}

func newInternalEvent(action internalEventActionType) internalEvent {
	return internalEvent{
		time: time.Now(),
	}
}

func (e internalEvent) Kind() string {
	return "internal"
}

func (e internalEvent) Action() string {
	return string(e.action)
}

func (e internalEvent) Related() any {
	return nil
}

func (e internalEvent) Time() time.Time {
	return e.time
}

var _ Event = internalEvent{}
