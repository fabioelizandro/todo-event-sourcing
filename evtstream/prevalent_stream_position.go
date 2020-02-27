package evtstream

type prevalentStreamPosition struct {
	value uint64
}

func newPrevalentStreamPosition(value uint64) *prevalentStreamPosition {
	return &prevalentStreamPosition{value: value}
}

func (i *prevalentStreamPosition) Value() interface{} {
	return i.value
}
