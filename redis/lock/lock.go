package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:       "0.0.0.0:6380",
		Password:   "",
		MaxRetries: 1,
		PoolSize:   2,
		DB:         2,
	})

	data := "bye"
	go func() {
		fmt.Println("setting 1...")
		defer fmt.Println("done 1...")
		v, err := client.SetNX("key", data+"1", 5*time.Second).Result()
		if err != nil {
			fmt.Println("failed", err)
		}
		fmt.Println("value 1: ", v)
	}()

	time.Sleep(2 * time.Second)
	fmt.Println()

	go func() {
		fmt.Println("setting 2...")
		defer fmt.Println("done 2...")
		v, err := client.SetNX("key", data+"2", 5*time.Second).Result()
		if err != nil {
			fmt.Println("failed", err)
		}
		fmt.Println("value 2: ", v)
	}()

	time.Sleep(1 * time.Second)

	fmt.Println("done: ", client.Get("key2"))

	select {
	case <-time.After(2 * time.Second):
		// end
	}
}
