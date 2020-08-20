package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	stan "github.com/nats-io/stan.go"
)

var natsURLs = []string{
	"nats://dev-gam-api.silkrode.com.tw:32002",
	//"nats://localhost:4222",
	//"nats://localhost:4223",
	// "nats://localhost:4224",
	// "nats://localhost:4225",
}

func main() {
	clusterID := "gam"
	clientID := uuid.New().String()
	topic := "queue-topic-1"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(strings.Join(natsURLs, ",")))
	if err != nil {
		panic(err)
	}
	defer sc.Close()

	var counter int
	for {
		// does not return until an ack has been received from NATS Streaming
		word := fmt.Sprintf("[%v] %d", time.Now(), counter)
		err := sc.Publish(topic, []byte(word))
		if err != nil {
			fmt.Println("Public error: ", err)

			sc.Close()
			sc, err = stan.Connect(clusterID, clientID, stan.NatsURL(strings.Join(natsURLs, ",")))
			if err != nil {
				panic(err)
			}
		}
		fmt.Println(word)

		time.Sleep(2 * time.Second)
		counter++
	}
}
