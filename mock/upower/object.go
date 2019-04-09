package upower

import (
	"github.com/godbus/dbus"
)

type DbusObject struct {
	IsCallStubCalled       bool
	CallStubArgs           CallStubArgs
	IsAddMatchSignalCalled bool
	AddMatchSignalStubArgs AddMatchSignalStubArgs
	ExpectedBrightess      int32
}

func NewDbusObject(expectedBrightness int32) *DbusObject {
	return &DbusObject{
		ExpectedBrightess: expectedBrightness,
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
	return &dbus.Call{
		Body: []interface{}{mobj.ExpectedBrightess},
		Err:  nil,
		Args: args,
	}
}
