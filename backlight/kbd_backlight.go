package backlight

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shadowbeetle/set-kbd-blight/upower"
	"github.com/godbus/dbus"
)

type KbdBacklight struct {
	*Config
	dbusObject        upower.DbusObject
	dbusSignalCh      chan *dbus.Signal
	desiredBrightness int32
	timer             *time.Timer
	inputCh           chan []byte
	ErrorCh           chan error
}

func NewKbdBacklight(conf Config) (*KbdBacklight, error) {
	err := conf.setDefaults()
	if err != nil {
		return nil, err
	}

	var initialBrightness int32
	brPtr := &initialBrightness
	busObject := upower.GetObject(conf.dbusConnection)
	call := upower.CallGetBrightness(busObject)
	err = upower.StoreBrightness(call, brPtr)
	if err != nil {
		return nil, err
	}

	dbusCh := make(chan *dbus.Signal, 10)
	upower.SignalListen(conf.dbusConnection, busObject, dbusCh)

	inputCh := make(chan []byte)
	errCh := make(chan error)

	kbl := &KbdBacklight{
		Config:            &conf,
		dbusObject:        busObject,
		dbusSignalCh:      dbusCh,
		desiredBrightness: initialBrightness,
		timer:             time.NewTimer(conf.IdleWaitTime),
		inputCh:           inputCh,
		ErrorCh:           errCh,
	}

	return kbl, nil
}

func (kbl *KbdBacklight) Run() error {
	var failCnt int
	for _, path := range kbl.InputPaths {
		f, err := os.Open(path)
		if err != nil {
			log.Println("could not open input", path, err.Error())
			failCnt += 1
			continue
		}

		go kbl.readInput(f)
	}

	if failCnt >= len(kbl.InputPaths) {
		return fmt.Errorf("could not open any of the provided inputs %v", kbl.InputPaths)
	}

	go kbl.listenUserBrightnessChange()
	go kbl.onIdleTurnOff()
	go kbl.onInputTurnOn()

	return nil
}

func (kbl *KbdBacklight) readInput(f *os.File) {
	b1 := make([]byte, 32)
	for {
		_, err := f.Read(b1)
		if err != nil {
			kbl.ErrorCh <- err
			continue
		}

		kbl.inputCh <- b1
	}
}

func (kbl *KbdBacklight) setBrightness() {
	upower.CallSetBrightness(kbl.dbusObject, kbl.desiredBrightness)
}

func (kbl *KbdBacklight) listenUserBrightnessChange() {
	for s := range kbl.dbusSignalCh {
		if s.Body[1] == "internal" {
			kbl.desiredBrightness = s.Body[0].(int32)
			kbl.timer.Reset(kbl.IdleWaitTime)
		}
	}
}

func (kbl *KbdBacklight) onIdleTurnOff() {
	for range kbl.timer.C {
		upower.CallSetBrightness(kbl.dbusObject, 0)
	}
}

func (kbl *KbdBacklight) onInputTurnOn() {
	for range kbl.inputCh {
		kbl.timer.Reset(kbl.IdleWaitTime)
		kbl.setBrightness()
	}
}
