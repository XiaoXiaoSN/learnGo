package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v8"
)

var mapKey = "myMapppp"

// LuaScript ...
var SetScript = redis.NewScript(`
local time_key = KEYS[1]
local map_key = KEYS[2]
local amount = tonumber(ARGV[1])

redis.call("INCRBY", time_key, amount)

local max_key = time_key .. "_mx"
local max = tonumber(redis.call("GET", max_key))
if max == nil or amount > max then
    redis.call("SET", max_key, amount)
end

redis.call("HSETNX", map_key, time_key, 1)

return ""
`)

var GetScript = redis.NewScript(`
local map_key = KEYS[1]
local cursor = "0"
local count = 1

local result = redis.call("HSCAN", map_key, cursor, "COUNT", count)
if result == nil then
    return "emp"
end

for _, value in pairs(result[2]) do 
    local amount = redis.call("GET", value)
    local max = redis.call("GET", value.."_mx")

    redis.call("DEL", value)
    redis.call("DEL", value.."_mx")
    redis.call("HDEL", map_key, value)

    return {value, amount, max}
end

return ""
`)

type trade struct {
	Amount  string
	TimeKey string
}

func main() {
	dataSource := make(chan trade, 100)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "10.1.1.111:6380",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	// get data from source
	go func() {
		for {
			dataSource <- trade{random(100), timeKey()}
			time.Sleep(time.Millisecond * 40)
		}
	}()

	// handle data
	go func() {
		var err error
		ctx := context.Background()

		for {
			select {
			case t := <-dataSource:
				_, err = SetScript.Run(ctx, redisClient, []string{t.TimeKey, mapKey}, t.Amount).Result()
				if err != nil && err != redis.Nil {
					fmt.Println(err)
				}
			}
		}
	}()

	// 取出資料
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			ctx := context.Background()

			result, err := GetScript.Run(ctx, redisClient, []string{mapKey}).Result()
			if err != nil && err != redis.Nil {
				fmt.Println(err)
			}
			fmt.Printf("result: %v\n", result)
		}
	}()

	select {}
}

func timeKey() string {
	t := int(time.Now().Unix())
	return strconv.Itoa(t - t%10)
}

func random(max int) string {
	return strconv.Itoa(rand.Int() % max)
}
