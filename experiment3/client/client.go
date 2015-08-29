package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

var url string

func main() {
	// parse commandl ine parameters
	flag.Parse()

	origin := "http://localhost/"
	url := flag.Arg(0)
	fmt.Printf("connecting to %q\n", url)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	var msg = make([]byte, 512)
	var n int
	for {
		if n, err = ws.Read(msg); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("Received: %s", msg[:n])
	}
}
