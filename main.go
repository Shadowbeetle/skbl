package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/godbus/dbus"
)

func main() {
	// Get brightness
	// dbus-send --type=method_call --print-reply=literal --system --dest="org.freedesktop.UPower" /org/freedesktop/UPower/KbdBacklight org.freedesktop.UPower.KbdBacklight.GetBrightness

	// Set birghtness
	//  dbus-send --type=method_call --print-reply=literal --system --dest="org.freedesktop.UPower" /org/freedesktop/UPower/KbdBacklight org.freedesktop.UPower.KbdBacklight.SetBrightness int32:2
	// Channel to hang in the end

	done := make(chan bool)
	defer func() { <-done }()

	check := func(e error) {
		if e != nil {
			close(done)
			panic(e)
		}
	}

	conn, err := dbus.SystemBus()
	check(err)

	initKbdBrightness := 0
	brPtr := &initKbdBrightness

	// Get initial keyboard backlight birghtness
	busObject := conn.Object("org.freedesktop.UPower", "/org/freedesktop/UPower/KbdBacklight")
	err = busObject.Call("org.freedesktop.UPower.KbdBacklight.GetBrightness", 0).Store(brPtr)
	check(err)

	// Set backlight brightness. Might need to make it a goroutine and use Object.Go instead. the value is between 0 and 3, now it's hardocded to set it to 3
	busObject.Call("org.freedesktop.UPower.KbdBacklight.SetBrightness", 0, int32(1))

	check(err)
	// Sender: :1.44, Path: /org/freedesktop/UPower/KbdBacklight, Name: org.freedesktop.UPower.KbdBacklight.BrightnessChangedWithSource, Body: [3 internal]

	// conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
	// "type='signal',path='/org/freedesktop/UPower/KbdBacklight',member='BrightnessChangedWithSource'")

	busObject.AddMatchSignal("org.freedesktop.UPower.KbdBacklight", "BrightnessChangedWithSource")
	dbusCh := make(chan *dbus.Signal, 10)
	conn.Signal(dbusCh)

	// handle SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Send bytes from mouse
	aggrCh := make(chan []byte)

	inputhPaths := []string{"/dev/input/mice", "/dev/input/event4"}

	for _, path := range inputhPaths {
		_, err := os.Stat(path)
		check(err)

		f, err := os.Open(path)
		check(err)
		go func() {
			for {
				b1 := make([]byte, 32) // Needs to be 32 long as the keyboard event is 32 bits
				_, err := f.Read(b1)
				check(err)
				aggrCh <- b1
			}
		}()
	}

	timer := time.NewTimer(5 * time.Second)
	// Process events
	go func() {
		for {
			select {
			case b := <-aggrCh:
				fmt.Printf("bytes: %s\n", string(b))
			case s := <-dbusCh:
				fmt.Printf("Body: %v, type of brightness: %T, type of origin: %T\n", s.Body, s.Body[0], s.Body[1])
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

	// Hang so doesn't exit
	// <-done
}

// TODO new flow: Provide input event path as comma ":" separated values and kbd_backglith event path relative to /dev/input/ in an ENV VAR. Add mice if you want to listen to all mice, or add the acutal mouse handler. For keyboards check /proc/bus/input/devices and look for kbd and check either "Handlers="and look for eventX.d). Listen for chages in those files, and send dbus event to set backlight, on change. Also set timer to turn it off after 5 seconds. Alos listen for backlight udev chnages as well, and update the desired brightness
