package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
	"golang.org/x/net/websocket"
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

	go buttonHandler(pin)

	// http listen
	//http.Handle("/button", websocket.Handler(EchoServer))
	http.HandleFunc("/button", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(EchoServer)}
		s.ServeHTTP(w, req)
	})
	http.ListenAndServe(":8080", nil)
}

// handle websocket connection
func EchoServer(ws *websocket.Conn) {
	fmt.Println("Client connected")

	w := bufio.NewWriter(ws)
	w.WriteString("Hello, i will tell you if button is pressed\n\n")
	w.Flush()

	for {
		<-ctrlChan
		w.WriteString("Kachhinggg... somebody pressed the button\n")
		w.Flush()
	}
}

// check if button is pressed
func buttonHandler(p gpio.Pin) {
	for {
		if p.Get() {
			ctrlChan <- true
		}
		time.Sleep(150 * time.Millisecond)
	}
}
