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
}

func NewTimer() *Timer {
	return &Timer{
		ResetStrobe:       make(chan bool),
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
