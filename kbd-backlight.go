package main

import (
	"log"
	"os"
	"time"

	"github.com/godbus/dbus"
)

type KbdBacklight struct {
	dbusObject        dbus.BusObject
	dbusSignalCh      chan *dbus.Signal
	desiredBrightness int32
	timer             *time.Timer
	inputCh           chan []byte
	idleWaitTime      time.Duration
	errorCh           chan error
}

func NewKbdBacklight(idleWaitTime time.Duration) (*KbdBacklight, error) {
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

	inputhPaths := []string{"/dev/input/mice", "/dev/input/event4"} // TODO make this settable by env vars or flags
	inputCh := make(chan []byte)
	errCh := make(chan error)

	kbr := &KbdBacklight{
		dbusObject:        busObject,
		dbusSignalCh:      dbusCh,
		desiredBrightness: initialBrightness,
		timer:             time.NewTimer(idleWaitTime),
		idleWaitTime:      idleWaitTime,
		inputCh:           inputCh,
		errorCh:           errCh,
	}

	for _, path := range inputhPaths {
		_, err := os.Stat(path)
		if err != nil {
			log.Println("Could not stat input", path, err)
			continue
		}

		f, err := os.Open(path)
		if err != nil {
			log.Println("Could not open input", path, err)
			continue
		}

		go kbr.readInput(f)
	}

	go kbr.listenUserBrightnessChange()
	go kbr.onIdleTurnOff()
	go kbr.onInputTurnOn()

	return kbr, nil
}

func (kbr *KbdBacklight) readInput(f *os.File) {
	b1 := make([]byte, 32) //TODO try to move this out of the loop // Needs to be 32 long as the keyboard event is 32 bits
	for {
		_, err := f.Read(b1)
		if err != nil {
			kbr.errorCh <- err
			continue
		}

		kbr.inputCh <- b1
	}
}

func (kbr *KbdBacklight) setBrightness() {
	kbr.dbusObject.Call("org.freedesktop.UPower.KbdBacklight.SetBrightness", 0, kbr.desiredBrightness)
}

func (kbr *KbdBacklight) listenUserBrightnessChange() {
	for s := range kbr.dbusSignalCh {
		if s.Body[1] == "internal" {
			kbr.desiredBrightness = s.Body[0].(int32)
			kbr.timer.Reset(kbr.idleWaitTime)
		}
	}
}

func (kbr *KbdBacklight) onIdleTurnOff() {
	for range kbr.timer.C {
		kbr.dbusObject.Call("org.freedesktop.UPower.KbdBacklight.SetBrightness", 0, 0)
	}
}

func (kbr *KbdBacklight) onInputTurnOn() {
	for range kbr.inputCh {
		kbr.timer.Reset(kbr.idleWaitTime)
		kbr.setBrightness()
	}
}
