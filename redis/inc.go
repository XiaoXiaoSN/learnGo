package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var pid string

func init() {
	flag.StringVar(&pid, "pid", "p0", "程式編號")
	flag.Parse()
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:       "0.0.0.0:6380",
		Password:   "",
		MaxRetries: 1,
		PoolSize:   2,
		DB:         1,
	})

	var data = "data-"

	go func() {
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
