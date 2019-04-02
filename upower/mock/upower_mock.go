package mock

import (
	"github.com/godbus/dbus"
)

type DbusConnection struct {
	IsObjectStubCalled bool
}

func NewDbusConnection() *DbusConnection {
	return &DbusConnection{
		IsObjectStubCalled: false,
	}
}

func (mconn *DbusConnection) Object(dest string, path dbus.ObjectPath) dbus.BusObject {
	mconn.IsObjectStubCalled = true
	return &dbus.Object{}
}

func (mconn *DbusConnection) Signal(ch chan<- *dbus.Signal) {}

func NewDbusObject(expectedBrightness int32) *DbusObject {
	return &DbusObject{
		ExpectedBrightess: expectedBrightness,
	}
}

type DbusObject struct {
	IsCallStubCalled       bool
	IsAddMatchSignalCalled bool
	AddMatchSignalStubArgs AddMatchSignalStubArgs
	ExpectedBrightess      int32
}

type AddMatchSignalStubArgs struct {
	Method string
	Member string
	Args   []dbus.MatchOption
}

func (mobj *DbusObject) AddMatchSignal(method string, member string, args ...dbus.MatchOption) *dbus.Call {
	mobj.AddMatchSignalStubArgs = AddMatchSignalStubArgs{method, member, args}
	return &dbus.Call{
		Body: []interface{}{},
		Err:  nil,
	}
}

func (mobj *DbusObject) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {
	return &dbus.Call{
		Body: []interface{}{mobj.ExpectedBrightess},
		Err:  nil,
		Args: args,
	}
}

// type MockDbusCall struct{}

// func (mcall *MockDbusCall) Store(retvalues ...interface{}) error {
// 	return errors.New("Error")
// }
