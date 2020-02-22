package eventstream

type inMemoryStreamReadResult struct {
	event        Event
	nextPosition StreamPosition
}

func (i *inMemoryStreamReadResult) Event() Event {
	return i.event
}

func (i *inMemoryStreamReadResult) NextStreamPosition() StreamPosition {
	return i.nextPosition
}
