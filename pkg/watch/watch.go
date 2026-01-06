package watch

type Interface interface {
	Stop()
	Watch() <-chan Event
}

var (
	DefaultChanSize int32 = 100
)

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
func (w emptyWatch) Watch() <-chan Event {
	return chan Event(w)
}
