package backlight

// import (
// 	"errors"

// 	"github.com/godbus/dbus"
// )

// type MockDbusConnection struct {
// 	isObjectStubCalled bool
// 	args               DbusObjectStubArgs
// }

// type DbusObjectStubArgs struct {
// 	dest string
// 	path dbus.ObjectPath
// }

// func (mconn *MockDbusConnection) Object(dest string, path dbus.ObjectPath) MockDbusObject {
// 	mconn.isObjectStubCalled = true
// 	return MockDbusObject{}
// }

// func (mconn *MockDbusConnection) Signal(ch chan<- *dbus.Signal) {}

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
