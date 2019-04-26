package backlight

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Shadowbeetle/set-kbd-blight/mock/clock"
	"github.com/Shadowbeetle/set-kbd-blight/mock/upower"
	"github.com/godbus/dbus"
)

func TestNewKbdBacklight(t *testing.T) {
	expBr := int32(1)
	expIdleWT := time.Duration(111) * time.Second

	expectedCallArgs := upower.CallStubArgs{
		Method: "org.freedesktop.UPower.KbdBacklight.GetBrightness",
	}

	expectedAddMatchSignalStubArgs := upower.AddMatchSignalStubArgs{
		Method: "org.freedesktop.UPower.KbdBacklight",
		Member: "BrightnessChangedWithSource",
	}

	mockConn := upower.NewDbusConnection()
	mockDObj := upower.NewDbusObject(expBr, true, false)

	conf := Config{
		IdleWaitTime:   expIdleWT,
		InputFiles:     []io.Reader{strings.NewReader("/test/input/kbd")},
		dbusConnection: mockConn,
		dbusObject:     mockDObj,
	}

	kbl, err := NewKbdBacklight(conf)

	if err != nil {
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

	if kbl.IdleWaitTime != expIdleWT {
		t.Errorf("expected kbl.IdleWaitTime to equl %v, got %v instead\n", expIdleWT, kbl.IdleWaitTime)
	}
}

func TestRunInput(t *testing.T) {
	done := make(chan bool)
	defer func() { done <- true }()

	expBr := int32(2)
	mockConn := upower.NewDbusConnection()
	mockDObj := upower.NewDbusObject(expBr, true, false)

	expectedCallArgs := upower.CallStubArgs{
		Method: "org.freedesktop.UPower.KbdBacklight.SetBrightness",
		Args:   []interface{}{expBr},
	}

	qwerInput := &strings.Reader{}
	asdfInput := &strings.Reader{}
	zxcvInput := &strings.Reader{}
	readers := []io.Reader{qwerInput, asdfInput, zxcvInput}
	timer := clock.NewTimer()

	conf := Config{
		IdleWaitTime:   time.Duration(222),
		InputFiles:     readers,
		dbusConnection: mockConn,
		dbusObject:     mockDObj,
		timer:          timer,
		timerC:         timer.C,
	}

	kbl, err := NewKbdBacklight(conf)

	if err != nil {
		t.Fatalf("expected nil error got %s instead\n", err.Error())
	}

	mockDObj.ShouldStore = false
	kbl.Run()

	go func() {
		for {
			select {
			case err := <-kbl.ErrorCh:
				if err == io.EOF {
					continue
				}
				t.Fatalf("got unexpected error from kbl.ErrorCh %s\n", err.Error())
			case <-done:
				return
			}
		}
	}()

	qwerInput.Reset("q")
	<-timer.ResetStrobe

	if timer.ResetStubArg != conf.IdleWaitTime {
		t.Errorf("expected timer.ResetStubArg to equal %v, got %v instead\n", conf.IdleWaitTime, timer.ResetStubArg)
	}

	if !reflect.DeepEqual(expectedCallArgs, mockDObj.CallStubArgs) {
		t.Errorf("expected mockDObj.CallStub to be called with %v, got %v instead", expectedCallArgs, mockDObj.CallStubArgs)
	}

	asdfInput.Reset("a")
	<-timer.ResetStrobe

	if timer.ResetStubArg != conf.IdleWaitTime {
		t.Errorf("expected timer.ResetStubArg to equal %v, got %v instead\n", conf.IdleWaitTime, timer.ResetStubArg)
	}

	if !reflect.DeepEqual(expectedCallArgs, mockDObj.CallStubArgs) {
		t.Errorf("expected mockDObj.CallStub to be called with %v, got %v instead", expectedCallArgs, mockDObj.CallStubArgs)
	}

	zxcvInput.Reset("z")
	<-timer.ResetStrobe

	if timer.ResetStubArg != conf.IdleWaitTime {
		t.Errorf("expected timer.ResetStubArg to equal %v, got %v instead\n", conf.IdleWaitTime, timer.ResetStubArg)
	}

	if !reflect.DeepEqual(expectedCallArgs, mockDObj.CallStubArgs) {
		t.Errorf("expected mockDObj.CallStub to be called with %v, got %v instead", expectedCallArgs, mockDObj.CallStubArgs)
	}

	qwerInput.Reset("w")
	<-timer.ResetStrobe

	if !reflect.DeepEqual(expectedCallArgs, mockDObj.CallStubArgs) {
		t.Errorf("expected mockDObj.CallStub to be called with %v, got %v instead", expectedCallArgs, mockDObj.CallStubArgs)
	}

	if timer.ResetStubArg != conf.IdleWaitTime {
		t.Errorf("expected timer.ResetStubArg to equal %v, got %v instead\n", conf.IdleWaitTime, timer.ResetStubArg)
	}

	if timer.ResetStubCallCount != 4 {
		t.Errorf("expected timer.Reset to be called 4 times got %d instead\n", timer.ResetStubCallCount)
	}

	if mockDObj.CallStubCallCount != 5 {
		t.Errorf("expected mockDObj.Call to be called 5 times got %d instead\n", mockDObj.CallStubCallCount)
	}
}

func TestRunIdle(t *testing.T) {
	done := make(chan bool)
	defer func() { done <- true }()

	idleWaitTime := time.Duration(333)
	mockConn := upower.NewDbusConnection()
	mockDObj := upower.NewDbusObject(3, true, false)

	expectedCallArgs := upower.CallStubArgs{
		Method: "org.freedesktop.UPower.KbdBacklight.SetBrightness",
		Args:   []interface{}{int32(0)},
	}

	input := &strings.Reader{}
	timer := clock.NewTimer()

	conf := Config{
		IdleWaitTime:   idleWaitTime,
		InputFiles:     []io.Reader{input},
		dbusConnection: mockConn,
		dbusObject:     mockDObj,
		timer:          timer,
		timerC:         timer.C,
	}

	kbl, err := NewKbdBacklight(conf)

	if err != nil {
		t.Fatalf("expected nil error got %s instead\n", err.Error())
	}

	mockDObj.ShouldStore = false
	mockDObj.ShouldCallStrobe = true
	kbl.Run()

	go func() {
		for {
			select {
			case err := <-kbl.ErrorCh:
				if err == io.EOF {
					continue
				}
				t.Fatalf("got unexpected error from kbl.ErrorCh %s\n", err.Error())
			case <-done:
				return
			}
		}
	}()

	timer.Expire()
	<-mockDObj.CallStrobe
	if !reflect.DeepEqual(expectedCallArgs, mockDObj.CallStubArgs) {
		t.Errorf("expected Call to be called with %+v got %+v instead", expectedCallArgs, mockDObj.CallStubArgs)
	}
}

func TestRunUserBrightnessChange(t *testing.T) {
	done := make(chan bool)
	defer func() { done <- true }()

	expectedBrightness := int32(10)
	idleWaitTime := time.Duration(444)
	mockConn := upower.NewDbusConnection()
	mockDObj := upower.NewDbusObject(4, true, false)

	input := &strings.Reader{}
	timer := clock.NewTimer()

	conf := Config{
		IdleWaitTime:   idleWaitTime,
		InputFiles:     []io.Reader{input},
		dbusConnection: mockConn,
		dbusObject:     mockDObj,
		timer:          timer,
		timerC:         timer.C,
	}

	kbl, err := NewKbdBacklight(conf)

	if err != nil {
		t.Fatalf("expected nil error got %s instead\n", err.Error())
	}

	mockDObj.ShouldStore = false
	mockDObj.ShouldCallStrobe = true
	kbl.Run()

	go func() {
		for {
			select {
			case err := <-kbl.ErrorCh:
				if err == io.EOF {
					continue
				}
				t.Fatalf("got unexpected error from kbl.ErrorCh %s\n", err.Error())
			case <-done:
				return
			}
		}
	}()

	kbl.dbusSignalCh <- &dbus.Signal{
		Body: []interface{}{expectedBrightness, "internal"},
	}

	<-timer.ResetStrobe

	if kbl.desiredBrightness != expectedBrightness {
		t.Errorf("expected kbl.desiredBrightness to be %d got %d instead", expectedBrightness, kbl.desiredBrightness)
	}
}

func TestConfig(t *testing.T) {}
