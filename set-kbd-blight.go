package main

import "log"

func main() {
	conf := ReadConfig()
	kbr, err := NewKbdBacklight(conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	for err = range kbr.errorCh {
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
