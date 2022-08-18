package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

//生产者消费者模式
func Producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		out <- factor * i
	}
}

func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	ch := make(chan int, 64)
	go Producer(3, ch)
	go Producer(5, ch)
	go Consumer(ch)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit %v\n", <-sig)
}
