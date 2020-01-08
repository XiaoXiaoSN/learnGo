package main

import (
	"fmt"

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
	err := client.Publish("fish", data).Err()
	if err != nil {
		// return errors.New("发布失败")
		fmt.Println("failed", err)
	}

	fmt.Println("done")
}
