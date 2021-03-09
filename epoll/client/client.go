package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8181", "http service address")
var count = 25000

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recover: ", debug.Stack(), e)
		}
	}()

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	conns := make([]*websocket.Conn, 0, count)
	for i := 0; i < count; i++ {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()
		conns = append(conns, c)

		if i%1000 == 0 {
			fmt.Println("now connected ", i)
		}
	}

	//
	fmt.Println("start send")

	done := make(chan struct{})
	for i := range conns {
		go func(c *websocket.Conn) {
			defer close(done)

			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recv: %s", message)
			}
		}(conns[i])
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// handle server response
	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:
			for i := range conns {
				go func(c *websocket.Conn) {
					err := c.WriteMessage(websocket.BinaryMessage, []byte(fmt.Sprintf("%v", time.Now().UnixNano()/1e6)))
					if err != nil {
						log.Println("write:", err)
					}
				}(conns[i])
			}

		case <-interrupt:
			log.Println("interrupt")

			for i := range conns {
				c := conns[i]

				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
