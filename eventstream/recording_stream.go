package eventstream

type recordingStream struct {
	stream EventStream
	memory []Event
}

func NewRecordingEventStream() *recordingStream {
	return &recordingStream{stream: NewPrevalentEventStream(
		newNoopPrevalentStreamStore(),
		[]*prevalentEventEnvelope{},
	)}
}

func (r *recordingStream) FirstPosition() StreamPosition {
	return r.stream.FirstPosition()
}

func (r *recordingStream) Read(position StreamPosition) (EventEnvelope, error) {
	return r.stream.Read(position)
}

func (r *recordingStream) ReadByCorrelationID(id string) ([]EventEnvelope, error) {
	return r.stream.ReadByCorrelationID(id)
}

func (r *recordingStream) Write(events []Event) error {
	err := r.stream.Write(events)
	if err != nil {
		return err
	}

	r.memory = append(r.memory, events...)

	return nil
}

func (r *recordingStream) Tape() []Event {
	return r.memory
}

func (r *recordingStream) EraseTape() {
	r.memory = []Event{}
}
