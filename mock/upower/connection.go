package upower

import "github.com/godbus/dbus"

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
