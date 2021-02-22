package main

import (
	"fmt"
	"sync"
	"time"
)

var c = SafeCounter{v: make(map[string]int)}
var theKey string = "key"

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	v   map[string]int
	mux sync.RWMutex
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mux.Unlock()
}

// Set ,,,
func (c *SafeCounter) Set(key string, val int) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key] = val
	c.mux.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mux.Unlock()
	return c.v[key]
}

// func run(i int) {
// 	c.mux.Lock()
// 	val, ok := c.v[theKey]
// 	c.mux.Unlock()

// 	if !ok {
// 		log.Printf("not ok: %d", i)
// 		c.mux.Lock()
// 		c.v[theKey] = 1
// 		c.mux.Unlock()
// 		log.Printf("write: %d", i)
// 		return
// 	}

// 	c.mux.Lock()
// 	c.v[theKey]++
// 	c.mux.Unlock()

// 	log.Printf("succ: %d %d", i, val)
// }

func run(i int) {
	c.mux.Lock()
	val, _ := c.v[theKey]
	if i == 0 {
		time.Sleep(time.Second)
	}
	c.mux.Unlock()

	c.Set(theKey, val+1)
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			run(i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Println(c.Value("somekey"))
}
