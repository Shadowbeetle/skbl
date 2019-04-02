package upower

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

func GetBrightness(o DbusObject) (int32, error) {
	var br int32
	err := o.Call(CALL_GET_BRIGHTNESS, 0).Store(&br)
	return br, err
}

func SetBrightness(o DbusObject, value int32) error {
	return o.Call(CALL_SET_BRIGHTNESS, 0, value).Store() // Hack so we don't need to listen on call.Dbus and to get call.Err returned instead of having it as a memeber
}

func SignalListen(conn DbusConnection, o DbusObject, ch chan<- *dbus.Signal) {
	o.AddMatchSignal(SIGNAL_INTERFACE, SIGNAL)
	conn.Signal(ch)
}
