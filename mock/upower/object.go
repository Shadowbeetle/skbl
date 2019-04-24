package upower

import (
	"github.com/godbus/dbus"
)

type DbusObject struct {
	IsCallStubCalled       bool
	CallStubArgs           CallStubArgs
	CallStubCallCount      int
	CallStrobe             chan bool
	ShouldCallStrobe       bool
	IsAddMatchSignalCalled bool
	AddMatchSignalStubArgs AddMatchSignalStubArgs
	ExpectedBrightess      int32
	ShouldStore            bool
}

func NewDbusObject(expectedBrightness int32, shouldStore bool, shouldCallStrobe bool) *DbusObject {
	return &DbusObject{
		ExpectedBrightess: expectedBrightness,
		ShouldStore:       shouldStore,
		CallStrobe:        make(chan bool),
		ShouldCallStrobe:  shouldCallStrobe,
	}
}

type AddMatchSignalStubArgs struct {
	Method string
	Member string
	Args   []dbus.MatchOption
}

func (mobj *DbusObject) AddMatchSignal(method string, member string, args ...dbus.MatchOption) *dbus.Call {
	mobj.IsAddMatchSignalCalled = true
	mobj.AddMatchSignalStubArgs = AddMatchSignalStubArgs{method, member, args}
	return &dbus.Call{
		Body: []interface{}{},
		Err:  nil,
	}
}

type CallStubArgs struct {
	Method string
	Flags  dbus.Flags
	Args   []interface{}
}

func (mobj *DbusObject) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {
	mobj.IsCallStubCalled = true
	mobj.CallStubArgs = CallStubArgs{method, flags, args}
	mobj.CallStubCallCount += 1

	if mobj.ShouldCallStrobe {
		mobj.CallStrobe <- true
	}

	var body []interface{}
	if mobj.ShouldStore {
		body = []interface{}{mobj.ExpectedBrightess}
	} else {
		body = []interface{}{}
	}

	return &dbus.Call{
		Body: body,
		Err:  nil,
		Args: args,
	}
}
