package main

import (
	"github.com/XiaoXiaoSN/learnGo/pprof/data"

	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		for {
			log.Println(data.Add("https://github.com/EDDYCJY"))
			time.Sleep(time.Second)
		}
	}()

	http.ListenAndServe("0.0.0.0:6060", nil)
}
