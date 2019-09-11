package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	data := make([]byte, 100000)
	var total int
	for {
		data = []byte(fmt.Sprintf("->%d ", total))
		n, err := conn.Write(data)
		if err != nil {
			total += n
			log.Printf("write %d bytes, error:%s\n", n, err)
			break
		}
		total += n
		log.Printf("write %d bytes this time, %d bytes in total\n\n", n, total)

		time.Sleep(time.Microsecond * 5)
	}

	log.Printf(">> write %d bytes in total\n", total)
	time.Sleep(time.Second * 10000)
}
