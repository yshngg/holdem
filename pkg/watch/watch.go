// xref: k8s.io/apimachinery/pkg/watch

package watch

type EventType int

const (
	internalRunFunctionMarker EventType = iota // internal do function
	Check
	Call
	Raise
	Fold
	AllIn
	Join
	Leave
	Start
	End
)

type Interface interface {
	Stop()
	ResultChan() <-chan Event
}

var (
	DefaultChanSize int32 = 100
)

type Event struct {
	Type   EventType
	Object any
}
type emptyWatch chan Event

func NewEmptyWatch() Interface {
	ch := make(chan Event)
	close(ch)
	return emptyWatch(ch)
}

// Stop implements Interface
func (w emptyWatch) Stop() {
}

// ResultChan implements Interface
func (w emptyWatch) ResultChan() <-chan Event {
	return chan Event(w)
}
