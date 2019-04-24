package clock

import (
	"time"
)

type Timer struct {
	IsResetStubCalled  bool
	ResetStubCallCount int
	ResetStubArg       time.Duration
	ResetStrobe        chan bool
	C                  <-chan time.Time
	Expire             func()
}

func NewTimer() *Timer {
	c := make(chan time.Time)

	return &Timer{
		ResetStrobe:       make(chan bool),
		C:                 c,
		Expire:            func() { c <- time.Time{} },
		IsResetStubCalled: false,
	}
}

func (t *Timer) Reset(d time.Duration) bool {
	t.IsResetStubCalled = true
	t.ResetStubCallCount += 1
	t.ResetStubArg = d
	t.ResetStrobe <- true
	return true
}
