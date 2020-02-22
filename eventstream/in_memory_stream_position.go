package eventstream

type inMemoryStreamPosition struct {
	value uint64
}

func (i *inMemoryStreamPosition) Value() interface{} {
	return i.value
}
