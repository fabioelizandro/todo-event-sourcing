package eventstream

type inMemoryStreamPosition struct {
	value uint64
}

func (i *inMemoryStreamPosition) Before(position StreamPosition) bool {
	value := position.Value().(uint64)
	return i.value < value
}

func (i *inMemoryStreamPosition) After(position StreamPosition) bool {
	value := position.Value().(uint64)
	return i.value > value
}

func (i *inMemoryStreamPosition) Value() interface{} {
	return i.value
}
