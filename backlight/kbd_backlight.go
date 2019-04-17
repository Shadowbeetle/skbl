package backlight

import (
	"fmt"
	"io"

	"github.com/Shadowbeetle/set-kbd-blight/upower"
	"github.com/godbus/dbus"
)

type KbdBacklight struct {
	*Config
	dbusSignalCh      chan *dbus.Signal
	desiredBrightness int32
	inputCh           chan bool
	ErrorCh           chan error
}

func NewKbdBacklight(conf Config) (*KbdBacklight, error) {
	err := conf.setDefaults()
	if err != nil {
		return nil, err
	}

	initBr, err := upower.GetBrightness(conf.dbusObject)
	if err != nil {
		return nil, err
	}

	dbusCh := make(chan *dbus.Signal, 10)
	upower.SignalListen(conf.dbusConnection, conf.dbusObject, dbusCh)

	inputCh := make(chan bool)
	errCh := make(chan error)

	kbl := &KbdBacklight{
		Config:            &conf,
		dbusSignalCh:      dbusCh,
		desiredBrightness: initBr,
		inputCh:           inputCh,
		ErrorCh:           errCh,
	}

	return kbl, nil
}

func (kbl *KbdBacklight) Run() error {
	for _, f := range kbl.InputFiles {
		go kbl.onInputTurnOn(f)
	}

	go kbl.onUserBrightnessChange()
	go kbl.onIdleTurnOff()

	return nil
}

func (kbl *KbdBacklight) onInputTurnOn(f io.Reader) {
	b1 := make([]byte, 32)
	for {
		_, err := f.Read(b1)
		if err != nil {
			kbl.ErrorCh <- err
			continue
		}

		fmt.Printf("%+v", kbl.timer)
		kbl.timer.Reset(kbl.IdleWaitTime)

		err = upower.SetBrightness(kbl.dbusObject, kbl.desiredBrightness)
		if err != nil {
			kbl.ErrorCh <- err
		}
	}
}

func (kbl *KbdBacklight) onIdleTurnOff() {
	for range kbl.timerC {
		err := upower.SetBrightness(kbl.dbusObject, 0)
		if err != nil {
			kbl.ErrorCh <- err
		}
	}
}

func (kbl *KbdBacklight) onUserBrightnessChange() {
	for s := range kbl.dbusSignalCh {
		if s.Body[1] == "internal" {
			kbl.desiredBrightness = s.Body[0].(int32)
			kbl.timer.Reset(kbl.IdleWaitTime)
		}
	}
}
