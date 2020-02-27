package evtstream

type prevalentStreamPosition struct {
	value int64
}

func newPrevalentStreamPosition(value int64) *prevalentStreamPosition {
	return &prevalentStreamPosition{value: value}
}

func (i *prevalentStreamPosition) After(position StreamPosition) bool {
	return i.value > position.Value()
}

func (i *prevalentStreamPosition) Before(position StreamPosition) bool {
	return i.value < position.Value()
}

func (i *prevalentStreamPosition) Next() StreamPosition {
	return newPrevalentStreamPosition(i.value + 1)
}

func (i *prevalentStreamPosition) Value() int64 {
	return i.value
}
