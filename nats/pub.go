package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	stan "github.com/nats-io/stan.go"
)

var counter int

func main() {
	clusterID := "test-cluster"
	clientID := uuid.New().String()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		fmt.Println(err)
	}
	// Close connection
	defer sc.Close()

	for {
		// does not return until an ack has been received from NATS Streaming
		word := fmt.Sprintf("[%v] %d", time.Now(), counter)
		sc.Publish("get.example.model", []byte(word))
		fmt.Println(word)

		time.Sleep(2 * time.Second)
		counter++
	}
}
