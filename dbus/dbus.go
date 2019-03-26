package dbus

import (
	"github.com/godbus/dbus"
)

const (
	DESTINATION         string          = "org.freedesktop.UPower"
	OBJECT_PATH         dbus.ObjectPath = "/org/freedesktop/UPower/KbdBacklight"
	CALL_GET_BRIGHTNESS string          = "org.freedesktop.UPower.KbdBacklight.GetBrightness"
	CALL_SET_BRIGHTNESS string          = "org.freedesktop.UPower.KbdBacklight.SetBrightness"
	SIGNAL_INTERFACE    string          = "org.freedesktop.UPower.KbdBacklight"
	SIGNAL              string          = "BrightnessChangedWithSource"
)

type DbusConnection interface {
	Object(string, dbus.ObjectPath) dbus.BusObject
	Signal(chan<- *dbus.Signal)
}

type DbusObject interface {
	Call(string, dbus.Flags, ...interface{}) *dbus.Call
	AddMatchSignal(string, string, ...dbus.MatchOption) *dbus.Call
}

type DbusCall interface {
	Store(...interface{}) error
}

func GetObject(conn DbusConnection) dbus.BusObject {
	return conn.Object(DESTINATION, OBJECT_PATH)
}

func CallGetBrightness(o DbusObject) *dbus.Call {
	return o.Call(CALL_GET_BRIGHTNESS, 0)
}

func CallSetBrightness(o DbusObject, value int32) {
	o.Call(CALL_SET_BRIGHTNESS, 0, value)
}

func StoreBrightness(c DbusCall, store *int32) error {
	return c.Store(store)
}

func SignalListen(conn DbusConnection, o DbusObject, ch chan<- *dbus.Signal) {
	o.AddMatchSignal(SIGNAL_INTERFACE, SIGNAL)
	conn.Signal(ch)
}
