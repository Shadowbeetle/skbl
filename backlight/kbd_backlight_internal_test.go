package backlight

import (
	"fmt"
	"testing"
	"time"

	"github.com/Shadowbeetle/set-kbd-blight/upower/mock"
	"github.com/godbus/dbus"
)

func TestNewKbdBacklight(t *testing.T) {
	// expectedConnDest := "org.freedesktop.UPower"
	expBr := int32(999)
	expectedAddMatchSignalStubArgs := mock.AddMatchSignalStubArgs{
		Method: "org.freedesktop.UPower.KbdBacklight",
		Member: "BrightnessChangedWithSource",
		Args:   []dbus.MatchOption{0},
	}

	mockConn := mock.NewDbusConnection()
	mockDObj := mock.NewDbusObject(expBr)

	conf := Config{
		IdleWaitTime:   time.Duration(5) * time.Second,
		InputPaths:     []string{"/test/input/kbd"},
		dbusConnection: mockConn,
		dbusObject:     mockDObj,
	}

	kbl, err := NewKbdBacklight(conf)
	if err != nil {
		panic(err)
		t.Fatalf("expected nil error got %s instead", err.Error())
	}

	if kbl.desiredBrightness != expBr {
		t.Errorf("expected kbl.desiredBrightess to equal %d got %d instead", expBr, kbl.desiredBrightness)
	}

	// TODO: should call busObject.AddMatchSignal("org.freedesktop.UPower.KbdBacklight", "BrightnessChangedWithSource")
	// TODO should call AddMatchSignal with expectedArgs
	fmt.Printf("%v\n", mockDObj)
	fmt.Printf("%v\n", kbl)
	// TODO: should retun &KbdBacklight with proper setup
}
