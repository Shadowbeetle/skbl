package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Shadowbeetle/set-kbd-blight/errorchecker"
	"github.com/godbus/dbus"
)

// TODO make this settable by env vars or flags
const idleWaitTime = 1 * time.Second

type KbdBrightness struct {
	dbusObject        dbus.BusObject
	dbusSignalCh      chan *dbus.Signal
	desiredBrightness int32
	timer             *time.Timer
	inputs            []os.File
	inputCh           chan []byte
	idleWaitTime      time.Duration
	errorChecker      errorchecker.ErrorChecker
}

func readInput(ec errorchecker.ErrorChecker, f *os.File, channel chan []byte) {
	for {
		b1 := make([]byte, 32) //TODO try to move this out of the loop // Needs to be 32 long as the keyboard event is 32 bits
		_, err := f.Read(b1)
		ec.Check(err)
		channel <- b1
	}
}

func NewKbdBrightness(ec errorchecker.ErrorChecker) *KbdBrightness {
	conn, err := dbus.SystemBus()
	ec.Check(err)

	var initialBrightness int32
	brPtr := &initialBrightness
	busObject := conn.Object("org.freedesktop.UPower", "/org/freedesktop/UPower/KbdBacklight")
	err = busObject.Call("org.freedesktop.UPower.KbdBacklight.GetBrightness", 0).Store(brPtr)
	ec.Check(err)

	busObject.AddMatchSignal("org.freedesktop.UPower.KbdBacklight", "BrightnessChangedWithSource")
	dbusCh := make(chan *dbus.Signal, 10)
	conn.Signal(dbusCh)

	inputhPaths := []string{"/dev/input/mice", "/dev/input/event4"} // TODO make this settable by env vars or flags
	inputCh := make(chan []byte)
	for _, path := range inputhPaths {
		_, err := os.Stat(path)
		ec.Check(err)

		f, err := os.Open(path)
		ec.Check(err)

		go readInput(ec, f, inputCh)
	}

	return &KbdBrightness{
		dbusObject:        busObject,
		dbusSignalCh:      dbusCh,
		desiredBrightness: initialBrightness,
		timer:             time.NewTimer(idleWaitTime),
		idleWaitTime:      idleWaitTime,
		inputCh:           inputCh,
		errorChecker:      ec,
	}
}

func (kbr *KbdBrightness) setBrightness() {
	kbr.dbusObject.Call("org.freedesktop.UPower.KbdBacklight.SetBrightness", 0, kbr.desiredBrightness)
}

func (kbr *KbdBrightness) listenUserBrightnessChange() {
	for s := range kbr.dbusSignalCh {
		if s.Body[1] == "internal" {
			kbr.desiredBrightness = s.Body[0].(int32)
			kbr.timer.Reset(kbr.idleWaitTime)
		}
	}
}

func (kbr *KbdBrightness) onIdleTurnOff() {
	for range kbr.timer.C {
		kbr.dbusObject.Call("org.freedesktop.UPower.KbdBacklight.SetBrightness", 0, 0)
	}
}

func (kbr *KbdBrightness) onInputTurnOn() {
	for range kbr.inputCh {
		kbr.timer.Reset(kbr.idleWaitTime)
		kbr.setBrightness()
	}
}

func main() {
	done := make(chan bool)
	defer func() { <-done }()

	ec := errorchecker.NewErrorChecker(done)

	kbr := NewKbdBrightness(ec)

	go kbr.listenUserBrightnessChange()
	go kbr.onIdleTurnOff()
	go kbr.onInputTurnOn()
	// conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
	// "type='signal',path='/org/freedesktop/UPower/KbdBacklight',member='BrightnessChangedWithSource'")

	// handle SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Send bytes from mouse

	timer := time.NewTimer(5 * time.Second)
	// Process events
	go func() {
		for {
			select {
			case <-timer.C:
				// fmt.Println("Timer elapsed, restarting")
				timer.Reset(5 * time.Second)
			case <-sigCh:
				fmt.Println("Receiveid SIGINT, halting")
				close(done)
				return
			}
		}
	}()
}
