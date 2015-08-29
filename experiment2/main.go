package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
)

// channel to control start blinking or not
var ctrlChan = make(chan bool)

func main() {
	// set GPIO25 to output mode
	pin, err := gpio.OpenPin(rpi.GPIO25, gpio.ModeOutput)
	if err != nil {
		fmt.Printf("Error opening pin! %s\n", err)
		return
	}

	// turn the led off on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Printf("\nClearing and unexporting the pin.\n")
			pin.Clear()
			pin.Close()
			os.Exit(0)
		}
	}()

	go dance(pin, ctrlChan)
	ctrlChan <- true

	// http listen
	http.HandleFunc("/dance", danceHandler)
	http.ListenAndServe(":8080", nil)
}

func dance(pin gpio.Pin, ctrlChan chan bool) {
	enabled := false
	for {
		select {
		case val := <-ctrlChan:
			fmt.Printf("dancing? %+v\n", val)
			enabled = val
		default:
			if enabled {
				pin.Set()
				time.Sleep(500 * time.Millisecond)
				pin.Clear()
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

func danceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request")
	s := r.URL.Query().Get("s")
	ctrlChan <- s == "1"
}
