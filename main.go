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

	kbl, err := backlight.NewKbdBacklight(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = kbl.Run()
	if err != nil {
		log.Fatal(err)
	}

	for err = range kbl.ErrorCh {
		if err != nil {
			log.Fatal(err)
		}
	}
}
