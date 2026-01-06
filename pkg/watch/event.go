package watch

import (
	"time"
)

type Event interface {
	Kind() string

	// Type() string

	// action is what action was taken/failed related to the related object.
	Action() string

	// regarding contains the object this Event is about.
	Related() any

	// time is the time when this Event was first observed
	Time() time.Time
}

type eventAction string

const (
	internalEventKind string = "internal" // internal event kind

	runFunction eventAction = "RunFunction" // internal do function\
)

type eventObject struct {
	do func()
}

type internalEvent struct {
	action eventAction
	object eventObject
	time   time.Time
}

func newInternalEvent(action eventAction, object eventObject) Event {
	return internalEvent{
		action: action,
		object: object,
		time:   time.Now(),
	}
}

func (e internalEvent) Kind() string {
	return internalEventKind
}

func (e internalEvent) Action() string {
	return string(e.action)
}

func (e internalEvent) Related() any {
	return e.object
}

func (e internalEvent) Time() time.Time {
	return e.time
}

var _ Event = internalEvent{}
