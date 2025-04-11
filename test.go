package main

import (
	"fmt"
	"time"
)

func producer(ch chan int) {
	for i := 1; ; i++ {
		ch <- i
		time.Sleep(time.Second)
	}
}

func consumer(ch chan int) {
	for {
		// 每消费 5 个数据时停止
		select {
		case v := <-ch:
			fmt.Println(v)
			if v == 5 {
				return
			}
		}
	}
}

//func main() {
//	ch := make(chan int, 10)
//
//	go producer(ch)
//	go consumer(ch)
//	select {}
//}
