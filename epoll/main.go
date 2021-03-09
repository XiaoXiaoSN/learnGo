package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var epoller *epoll

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// for {
	// 	echo(conn)
	// }
	if err := epoller.Add(conn); err != nil {
		log.Printf("Failed to add connection")
		conn.Close()
	}
}

func main() {
	// Increase resources limitations
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6061", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	// Start epoll
	var err error
	epoller, err = MkEpoll()
	if err != nil {
		panic(err)
	}

	go start()

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println("max connections:", maxConn)
		}
	}()

	http.HandleFunc("/echo", wsHandler)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}
}

func start() {
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}
		for i, conn := range connections {
			if conn == nil {
				break
			}
			go echo(connections[i])
		}
	}
}

func echo(conn *websocket.Conn) {
	code, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("[go] 開始 remove 囉~~~~~~~~~~")
		if err := epoller.Remove(conn); err != nil {
			log.Printf("Failed to remove %v", err)
		}
		conn.Close()
	} else {
		// log.Printf("msg: %s", string(msg))
		sentTime, _ := strconv.ParseInt(string(msg), 10, 64)
		usedTime := time.Now().UnixNano()/1e6 - sentTime
		a := new(bytes.Buffer)
		_ = binary.Write(a, binary.LittleEndian, usedTime)
		msg = a.Bytes()
	}

	err = conn.WriteMessage(code, msg)
	if err != nil {
		log.Printf("conn.WriteMessage msg")
	}
}
