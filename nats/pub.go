package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

var natsURLs = []string{
	"nats://localhost:4223",
	// "nats://localhost:4224",
	// "nats://localhost:4225",
}

func main() {
	topic := "topic.3.acc.1"

	sc, err := nats.Connect(strings.Join(natsURLs, ","))
	if err != nil {
		fmt.Println("nats.Connect >  ", err)
		return
	}
	defer sc.Close()

	var counter int
	for {
		// does not return until an ack has been received from NATS Streaming
		word := fmt.Sprintf("[%v] %d", time.Now(), counter)
		err := sc.Publish(topic, []byte(word))
		if err != nil {
			fmt.Println("Public error: ", err)
		}
		fmt.Println(word)

		time.Sleep(2 * time.Second)
		counter++
	}
}
