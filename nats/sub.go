package main

import (
	"fmt"

	"github.com/google/uuid"
	stan "github.com/nats-io/stan.go"
)

func main() {
	clusterID := "test-cluster"
	clientID := uuid.New().String()
	clientID = "ddd"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		fmt.Println("stan.Connect >  ", err)
		return
	}
	// Close connection
	defer sc.Close()

	// td, _ := time.ParseDuration("10s")

	// Simple Async Subscriber
	sub, _ := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	}, stan.DurableName("my-durable"))

	// Unsubscribe
	defer sub.Unsubscribe()

	// hold on
	select {}
}
