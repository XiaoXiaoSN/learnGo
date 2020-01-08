package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// SubscribeChannel redis topic
func SubscribeChannel(redisClient *redis.Client) {
	pubsub := redisClient.Subscribe("fish")
	_, err := pubsub.Receive()
	if err != nil {
		fmt.Println(err)
		return
	}

	timeout := time.After(10 * time.Second)
	ch := pubsub.Channel()
	x := make(chan string)
	for {
		select {
		case x3 := <-x:
			fmt.Println(x3)
		case msg := <-ch:
			fmt.Println(msg.Channel, msg.Payload, "\r\n")
			timeout = time.After(5 * time.Second)
		case <-timeout:
			fmt.Println("timeout")
			return
		}
	}

	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload, "\r\n")
		if msg.Payload == "bye" {
			fmt.Println("bye bye")
			return
		}
	}
}

func main() {
	fmt.Println("init")

	client := redis.NewClient(&redis.Options{
		Addr:       "0.0.0.0:6380",
		Password:   "",
		MaxRetries: 1,
		PoolSize:   2,
		DB:         2,
	})

	SubscribeChannel(client)

	fmt.Println("done")
}
