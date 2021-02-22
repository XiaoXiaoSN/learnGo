// https://stackoverflow.com/questions/47445344/is-there-a-difference-in-go-between-a-counter-using-atomic-operations-and-one-us

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Counter interface {
	Inc()
	Load() int64
}

// Atomic Implementation

type AtomicCounter struct {
	counter int64
}

func (c *AtomicCounter) Inc() {
	atomic.AddInt64(&c.counter, 1)
}

func (c *AtomicCounter) Load() int64 {
	return atomic.LoadInt64(&c.counter)
}

// Mutex Implementation

type MutexCounter struct {
	counter int64
	lock    sync.Mutex
}

func (c *MutexCounter) Inc() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.counter++
}

func (c *MutexCounter) Load() int64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.counter
}

// Use them...

func main() {
	fmt.Println("Running...")

	fmt.Println("Atomic Counter with Sleep")
	TestCounter_WithSleep(new(AtomicCounter))

	fmt.Println("Atomic Counter with Wait Group")
	TestCounter_WithWaitGroup(new(AtomicCounter))

	fmt.Println("Mutex Counter with Sleep")
	TestCounter_WithSleep(new(MutexCounter))

	fmt.Println("Mutex Counter with Wait Group")
	TestCounter_WithWaitGroup(new(MutexCounter))

	fmt.Println("Atomic Counter Print to 50")
	TestCounter_PrintAll(new(AtomicCounter))
	time.Sleep(time.Second)

	fmt.Println("Mutex Counter Print to 50")
	TestCounter_PrintAll(new(MutexCounter))
	time.Sleep(time.Second)

	fmt.Println("Atomic Counter Print to 50 in Go Routine")
	TestCounter_PrintAllGoRoutine(new(AtomicCounter))
	time.Sleep(time.Second)

	fmt.Println("Mutex Counter Print to 50 in Go Routine")
	TestCounter_PrintAllGoRoutine(new(MutexCounter))
	time.Sleep(time.Second)
}

func TestCounter_WithSleep(counter Counter) {
	for i := 0; i < 50; i++ {

		go func() {
			counter.Inc()
		}()
	}

	time.Sleep(time.Second)

	fmt.Println("Count:", counter.Load())
}

func TestCounter_WithWaitGroup(counter Counter) {
	wg := new(sync.WaitGroup)

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			counter.Inc()
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("Count:", counter.Load())
}

func TestCounter_PrintAll(counter Counter) {
	for i := 0; i < 50; i++ {
		go func() {
			counter.Inc()
			fmt.Println(counter.Load())
		}()
	}
}

func TestCounter_PrintAllGoRoutine(counter Counter) {
	for i := 0; i < 50; i++ {
		go func() {
			counter.Inc()
			go fmt.Println(counter.Load())
		}()
	}
}
