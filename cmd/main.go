package main

import (
	"log"

	"github.com/Shadowbeetle/set-kbd-blight/backlight"
)

func main() {
	conf := &backlight.Config{
		InputPaths:   inputPaths,
		IdleWaitTime: idleWaitTime,
	}

	kbr, err := backlight.NewKbdBacklight(conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	for err = range kbr.errorCh {
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
