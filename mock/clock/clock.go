package clock

import (
	"time"
)

type Timer struct {
	IsResetStubCalled  bool
	ResetStubCallTimes int
	ResetStubArg       time.Duration
	ResetStrobe        chan bool
}

func NewTimer() *Timer {
	return &Timer{
		ResetStrobe:       make(chan bool),
		IsResetStubCalled: false,
	}
}

func (t *Timer) Reset(d time.Duration) bool {
	t.IsResetStubCalled = true
	t.ResetStubCallTimes += 1
	t.ResetStubArg = d
	t.ResetStrobe <- true
	return true
}
