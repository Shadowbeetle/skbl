package main

func main() {
	conf := ReadConfig()
	kbr, err := NewKbdBacklight(conf)
	if err != nil {
		panic(err)
	}

	for err = range kbr.errorCh {
		if err != nil {
			panic(err)
		}
	}
}
