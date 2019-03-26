package backlight

import (
	"time"

	"github.com/Shadowbeetle/set-kbd-blight/dbus"
	godbus "github.com/godbus/dbus"
)

type Config struct {
	IdleWaitTime   time.Duration
	InputPaths     []string
	dbusConnection dbus.DbusConnection
}

func (conf *Config) setDefaults() error {
	if conf.dbusConnection == nil {
		dConn, err := godbus.SystemBus()
		if err != nil {
			return err
		}
		conf.dbusConnection = dConn
	}
	return nil
}
