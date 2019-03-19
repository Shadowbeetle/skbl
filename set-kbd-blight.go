package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// TODO make this settable by env vars or flags
const idleWaitTime = 1 * time.Second

func main() {
	viper := viper.New()
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.SetEnvPrefix("KBDBL")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	waitSeconds := time.Duration(viper.GetInt("wait-seconds")) * time.Second
	inputs := viper.GetStringSlice("inputs") // KBDBL_INPUTS=comma,separated,values

	fmt.Printf("inputs", inputs, "waitSeconds", waitSeconds, "env", viper.Get("wait_seconds"))
	// ------ READ CONFIG -------

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
