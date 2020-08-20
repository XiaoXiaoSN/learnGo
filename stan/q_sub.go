package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	stan "github.com/nats-io/stan.go"
)

var natsURLs = []string{
	"nats://dev-gam-api.silkrode.com.tw:32002",
	// "nats://localhost:4223",
	// "nats://localhost:4224",
	// "nats://localhost:4225",
}

func main() {
	clusterID := "gam"
	clientID := uuid.New().String()
	topic := "queue-topic-1"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(strings.Join(natsURLs, ",")))
	if err != nil {
		fmt.Println("stan.Connect >  ", err)
		return
	}
	defer sc.Close()

	// queue 1
	// Create a queue subscriber on "foo" for group "group1"
	go func() {
		queueID := 1
		// qsub, _ := sc.Subscribe(topic, func(m *stan.Msg) {
		qsub, _ := sc.QueueSubscribe(topic, "group1", func(m *stan.Msg) {
			fmt.Printf("[%d] Received a message: %s\n", queueID, string(m.Data))
		})
		defer qsub.Unsubscribe()

		// hold on gorutine
		select {}
	}()

	// queue 2
	// Create a queue subscriber on "foo" for group "group1"
	go func() {
		queueID := 2
		// qsub, _ := sc.Subscribe(topic, func(m *stan.Msg) {
		qsub, _ := sc.QueueSubscribe(topic, "group2", func(m *stan.Msg) {
			fmt.Printf("[%d] Received a message: %s\n", queueID, string(m.Data))
		})
		defer qsub.Unsubscribe()

		// hold on gorutine
		select {}
	}()

	// hold on service
	select {}
}
