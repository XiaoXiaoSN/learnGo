package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	stan "github.com/nats-io/stan.go"
)

var natsURLs = []string{
	//	"nats://dev-gam-api.silkrode.com.tw:32002",
	"nats://localhost:4222",
	//"nats://localhost:4223",
	// "nats://localhost:4224",
	// "nats://localhost:4225",
}

func main() {
	clusterID := "gam"
	clientID := uuid.New().String()
	topic := "@@pop.topics"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(strings.Join(natsURLs, ",")))
	if err != nil {
		fmt.Println("stan.Connect err: ", err)
		return
	}
	defer sc.Close()

	count := 0
	// Simple Async Subscriber
	sub, err := sc.Subscribe(topic, func(m *stan.Msg) {
		count++
		fmt.Printf("Received message count: %d\n", count)
	}, stan.DurableName("demoClient"), stan.StartAtSequence(0))
	if err != nil {
		fmt.Println("Subscribe error: ", err)
		return
	}
	defer sub.Unsubscribe()

	fmt.Println("start to subscribe...")

	// hold on
	select {}
}
