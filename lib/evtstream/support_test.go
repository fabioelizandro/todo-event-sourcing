package evtstream_test

type SomethingHappened struct {
	ID   string
	Data string
}

func (s *SomethingHappened) Type() string {
	return "SOMETHING_HAPPENED"
}

func (s *SomethingHappened) CorrelationID() string {
	return s.ID
}

type SomethingElseHappened struct {
	ID   string
	Data string
}

func (s *SomethingElseHappened) Type() string {
	return "SOMETHING_ELSE_HAPPENED"
}

func (s *SomethingElseHappened) CorrelationID() string {
	return s.ID
}
