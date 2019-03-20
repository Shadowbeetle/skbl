package backlight

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/godbus/dbus"
)

type Config struct {
	IdleWaitTime time.Duration
	InputPaths   []string
}

type KbdBacklight struct {
	dbusObject        dbus.BusObject
	dbusSignalCh      chan *dbus.Signal
	desiredBrightness int32
	timer             *time.Timer
	inputCh           chan []byte
	idleWaitTime      time.Duration
	errorCh           chan error
}

func NewKbdBacklight(conf *Config) (*KbdBacklight, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	var initialBrightness int32
	brPtr := &initialBrightness
	busObject := conn.Object("org.freedesktop.UPower", "/org/freedesktop/UPower/KbdBacklight")
	err = busObject.Call("org.freedesktop.UPower.KbdBacklight.GetBrightness", 0).Store(brPtr)
	if err != nil {
		return nil, err
	}

	busObject.AddMatchSignal("org.freedesktop.UPower.KbdBacklight", "BrightnessChangedWithSource")
	dbusCh := make(chan *dbus.Signal, 10)
	conn.Signal(dbusCh)

	inputCh := make(chan []byte)
	errCh := make(chan error)

	kbl := &KbdBacklight{
		dbusObject:        busObject,
		dbusSignalCh:      dbusCh,
		desiredBrightness: initialBrightness,
		timer:             time.NewTimer(conf.IdleWaitTime),
		idleWaitTime:      conf.IdleWaitTime,
		inputCh:           inputCh,
		errorCh:           errCh,
	}

	var failCnt int
	for _, path := range conf.InputPaths {
		_, err := os.Stat(path)
		if err != nil {
			log.Println("Could not stat input", path, err.Error())
			failCnt += 1
			continue
		}

		f, err := os.Open(path)
		if err != nil {
			log.Println("Could not open input", path, err.Error())
			failCnt += 1
			continue
		}

		go kbl.readInput(f)
	}

	if failCnt >= len(conf.InputPaths) {
		return nil, fmt.Errorf("Could not open any of the provided inputs %v", conf.InputPaths)
	}

	go kbl.listenUserBrightnessChange()
	go kbl.onIdleTurnOff()
	go kbl.onInputTurnOn()

	return kbl, nil
}

func (kbl *KbdBacklight) readInput(f *os.File) {
	b1 := make([]byte, 32)
	for {
		_, err := f.Read(b1)
		if err != nil {
			kbl.errorCh <- err
			continue
		}

		kbl.inputCh <- b1
	}
}

func (kbl *KbdBacklight) setBrightness() {
	kbl.dbusObject.Call("org.freedesktop.UPower.KbdBacklight.SetBrightness", 0, kbl.desiredBrightness)
}

func (kbl *KbdBacklight) listenUserBrightnessChange() {
	for s := range kbl.dbusSignalCh {
		if s.Body[1] == "internal" {
			kbl.desiredBrightness = s.Body[0].(int32)
			kbl.timer.Reset(kbl.idleWaitTime)
		}
	}
}

func (kbl *KbdBacklight) onIdleTurnOff() {
	for range kbl.timer.C {
		kbl.dbusObject.Call("org.freedesktop.UPower.KbdBacklight.SetBrightness", 0, 0)
	}
}

func (kbl *KbdBacklight) onInputTurnOn() {
	for range kbl.inputCh {
		kbl.timer.Reset(kbl.idleWaitTime)
		kbl.setBrightness()
	}
}
