package main

import (
	"io"
	"log"
	"os"

	"github.com/Shadowbeetle/set-kbd-blight/backlight"
)

func main() {
	var failCnt int
	var inputFiles []io.Reader
	for _, path := range inputPaths {
		f, err := os.Open(path)
		if err != nil {
			log.Println("could not open input", path, err.Error())
			failCnt += 1
			continue
		}
		inputFiles = append(inputFiles, f)
	}

	if failCnt >= len(inputFiles) {
		log.Fatalf("could not open any of the provided inputs %v", inputPaths)
	}

	conf := backlight.Config{
		InputFiles:   inputFiles,
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
