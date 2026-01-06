package watch

import "sync"

type Recorder interface {
	Events() []Event

	Interface
}

// Recorder records all events that are sent from the watch until it is closed.
type recorder struct {
	Interface

	lock   sync.Mutex
	events []Event
}

var _ Interface = &recorder{}

// NewRecorder wraps an Interface and records any changes sent across it.
func NewRecorder(w Interface) Recorder {
	r := &recorder{}
	r.Interface = Filter(w, r.record)
	return r
}

// record is a FilterFunc and tracks each received event.
func (r *recorder) record(in Event) (Event, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.events = append(r.events, in)
	return in, true
}

// Events returns a copy of the events sent across this recorder.
func (r *recorder) Events() []Event {
	r.lock.Lock()
	defer r.lock.Unlock()
	copied := make([]Event, len(r.events))
	copy(copied, r.events)
	return copied
}
