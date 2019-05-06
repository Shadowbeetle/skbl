package backlight

import (
	"io"
	"time"

	"github.com/Shadowbeetle/skbl/clock"
	"github.com/Shadowbeetle/skbl/upower"
	"github.com/godbus/dbus"
)

type Config struct {
	IdleWaitTime   time.Duration
	InputFiles     []io.Reader
	dbusConnection upower.DbusConnection
	dbusObject     upower.DbusObject
	timer          clock.Timer
	timerC         <-chan time.Time
}

type preventListeners struct {
	userBrightnessChange bool
	idle                 bool
	input                bool
}

func (conf *Config) setDefaults() error {
	if conf.dbusConnection == nil {
		dConn, err := dbus.SystemBus()
		if err != nil {
			return err
		}
		conf.dbusConnection = dConn
	}

	if conf.dbusObject == nil {
		conf.dbusObject = upower.GetObject(conf.dbusConnection)
	}

	var timer *time.Timer

	if conf.timer == nil && conf.timerC != nil {
		return TimerConfigError
	}

	if conf.timer == nil {
		timer = time.NewTimer(conf.IdleWaitTime)
		conf.timer = timer
		conf.timerC = timer.C
	}

	return nil
}
