package backlight

import (
	"time"

	"github.com/Shadowbeetle/set-kbd-blight/upower"
	"github.com/godbus/dbus"
)

type Config struct {
	IdleWaitTime   time.Duration
	InputPaths     []string
	dbusConnection upower.DbusConnection
	dbusObject     upower.DbusObject
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

	return nil
}
