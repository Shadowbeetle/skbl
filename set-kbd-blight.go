package main

import (
	"os"
	"os/signal"
	"time"
)

// TODO make this settable by env vars or flags
const idleWaitTime = 1 * time.Second

func main() {
	kbr, err := NewKbdBacklight(idleWaitTime)
	if err != nil {
		panic(err)
	}

	// handle SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	for {
		select {
		case err := <-kbr.errorCh:
			panic(err)
		case <-sigCh:
			return
		}
	}
}
