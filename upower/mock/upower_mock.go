package mock

import (
	"github.com/godbus/dbus"
)

type DbusConnection struct {
	IsObjectStubCalled bool
	Args               DbusObjectStubArgs
}

type DbusObjectStubArgs struct {
	Dest string
	Path dbus.ObjectPath
}

func NewDbusConnection() *DbusConnection {
	return &DbusConnection{
		IsObjectStubCalled: false,
		Args:               DbusObjectStubArgs{},
	}
}

func (mconn *DbusConnection) Object(dest string, path dbus.ObjectPath) dbus.Object {
	mconn.IsObjectStubCalled = true
	mconn.Args = DbusObjectStubArgs{dest, path}
	return dbus.Object{}
}

func (mconn *DbusConnection) Signal(ch chan<- *dbus.Signal) {}

// type MockDbusObject struct {
// 	isCallStubCalled bool
// }

// func (mobj *MockDbusObject) AddMatchSignal(iface string, member string) *MockDbusCall {
// 	return &MockDbusCall{}
// }

// func (movj *MockDbusObject) Call(method string, flags dbus.Flags, args ...interface{}) *MockDbusCall {
// 	return &MockDbusCall{}
// }

// type MockDbusCall struct{}

// func (mcall *MockDbusCall) Store(retvalues ...interface{}) error {
// 	return errors.New("Error")
// }
