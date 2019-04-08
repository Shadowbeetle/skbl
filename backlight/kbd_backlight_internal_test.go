package backlight

import (
	"reflect"
	"testing"
	"time"

	"github.com/Shadowbeetle/set-kbd-blight/upower/mock"
)

func TestNewKbdBacklight(t *testing.T) {
	expBr := int32(999)
	expIdleWT := time.Duration(5) * time.Second

	expectedCallArgs := mock.CallStubArgs{
		Method: "org.freedesktop.UPower.KbdBacklight.GetBrightness",
	}

	expectedAddMatchSignalStubArgs := mock.AddMatchSignalStubArgs{
		Method: "org.freedesktop.UPower.KbdBacklight",
		Member: "BrightnessChangedWithSource",
	}

	mockConn := mock.NewDbusConnection()
	mockDObj := mock.NewDbusObject(expBr)

	conf := Config{
		IdleWaitTime:   expIdleWT,
		InputPaths:     []string{"/test/input/kbd"},
		dbusConnection: mockConn,
		dbusObject:     mockDObj,
	}

	kbl, err := NewKbdBacklight(conf)

	if err != nil {
		panic(err)
		t.Fatalf("expected nil error got %s instead\n", err.Error())
	}

	if !mockDObj.IsCallStubCalled {
		t.Fatalf("expected Call to be called\n")
	}

	if !reflect.DeepEqual(expectedCallArgs, mockDObj.CallStubArgs) {
		t.Fatalf("expected Call to be called with %v got %v instead\n", expectedCallArgs, mockDObj.CallStubArgs)
	}

	if kbl.desiredBrightness != expBr {
		t.Errorf("expected kbl.desiredBrightess to equal %d got %d instead\n", expBr, kbl.desiredBrightness)
	}

	if !mockDObj.IsAddMatchSignalCalled {
		t.Fatalf("expeceted AddMatchSignal to be called\n")
	}

	if !reflect.DeepEqual(expectedAddMatchSignalStubArgs, mockDObj.AddMatchSignalStubArgs) {
		t.Errorf("expected AddMatchSignal to be called with %v got %v instead\n", expectedAddMatchSignalStubArgs, mockDObj.AddMatchSignalStubArgs)
	}

	if reflect.DeepEqual(kbl.Config, conf) {
		t.Errorf("expected kbl.Config to be %+v got %+v instead\n", kbl.Config, conf)
	}

	if kbl.inputCh == nil {
		t.Errorf("expected kbl.inputCh to be set, got nil instead\n")
	}

	if kbl.dbusSignalCh == nil {
		t.Errorf("expected kbl.dbusSignalCh to be set, got nil instead\n")
	}

	if kbl.timer == nil {
		t.Errorf("expected kbl.timer to be set, got nil isntead\n")
	}

	if kbl.ErrorCh == nil {
		t.Errorf("expected kbl.ErrorCh to be set, got nil isntead\n")
	}

	if kbl.IdleWaitTime != expIdleWT {
		t.Errorf("expected kbl.IdleWaitTime to equl %v, got %v instead\n", expIdleWT, kbl.IdleWaitTime)
	}
}
