package evtstream

import "time"

type StreamClock interface {
	Now() time.Time
}

type utcCLock struct {
}

func NewUTCCLock() *utcCLock {
	return &utcCLock{}
}

func (u *utcCLock) Now() time.Time {
	return time.Now().UTC()
}

type frozenClock struct {
	clockTime time.Time
}

func NewFrozenClock(clockTime time.Time) *frozenClock {
	return &frozenClock{clockTime: clockTime}
}

func (f *frozenClock) Now() time.Time {
	return f.clockTime
}
