package main

import (
	"fmt"
	"time"

	"github.com/yanyiwu/gojieba"
)

var jieba = gojieba.NewJieba()

func main() {
	tm := time.NewTimer(time.Second * 10)

	for {
		select {
		case <-tm.C:
			return
		default:
			go jiebaaaaa()
		}
		time.Sleep(time.Millisecond * 100) // 0.1s
	}
}

func jiebaaaaa() {
	// var jieba = gojieba.NewJieba()

	useHMM := true
	s := jieba.Cut("对应支付类别的帐号", useHMM)
	fmt.Println(s)
}
