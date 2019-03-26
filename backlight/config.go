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
}

func (conf *Config) setDefaults() error {
	if conf.dbusConnection == nil {
		dConn, err := dbus.SystemBus()
		if err != nil {
			return err
		}
		conf.dbusConnection = dConn
	}
	return nil
}
