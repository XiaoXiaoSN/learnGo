package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nats-io/nats.go"
)

var natsURLs = []string{
	"nats://localhost:4223",
	// "nats://localhost:4224",
	// "nats://localhost:4225",
}

func main() {
	topic := "@@pop.>"

	sc, err := nats.Connect(strings.Join(natsURLs, ","))
	if err != nil {
		fmt.Println("nats.Connect >  ", err)
		return
	}
	defer sc.Close()

	// Simple Async Subscriber
	sub, err := sc.Subscribe(topic, func(m *nats.Msg) {
		d, _ := DeflateDecode(m.Data)
		fmt.Printf("Received a message: %s\n", string(d))
		m.Respond([]byte("got you"))
	})
	if err != nil {
		fmt.Println("Subscribe error: ", err)
	}

	// Unsubscribe
	defer sub.Unsubscribe()

	fmt.Println("subscribe ok...")

	// hold on
	select {}
}

func DeflateDecode(input []byte) (result []byte, err error) {
	return ioutil.ReadAll(flate.NewReader(bytes.NewReader(input)))
}
