package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		fmt.Println("Start to watring connection")
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}

		// start a new goroutine to handle the new connection.
		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	// time.Sleep(time.Second * 10)
	for {
		// read from the connection
		var buf = make([]byte, 2000000)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
		}

		log.Printf("read %d bytes, content is %s\n\n", n, string(buf[:20]))
		time.Sleep(20 * time.Second)
	}
}
