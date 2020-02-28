package evtstream

import "time"

type recordingStream struct {
	stream EventStream
	tape   []Event
}

func NewRecordingEventStream() *recordingStream {
	return &recordingStream{
		stream: NewPrevalentEventStream(newNoopPrevalentStreamStore(),
			[]*prevalentEventEnvelope{},
			NewFrozenClock(time.Unix(946684800, 0).UTC()),
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

	r.tape = append(r.tape, events...)

	return nil
}

func (r *recordingStream) Tape() []Event {
	return r.tape
}

func (r *recordingStream) EraseTape() {
	r.tape = []Event{}
}
