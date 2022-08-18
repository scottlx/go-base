package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg = sync.WaitGroup{}

func busi(ch chan int, i int) {

	for i := range ch {
		fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())
		wg.Done()
	}

}

func startTask(ch chan int, i int) {
	wg.Add(1)
	ch <- i
}
func main() {

	task_cnt := 10
	go_cnt := 3
	ch := make(chan int)

	//工作协程池
	for i := 0; i < go_cnt; i++ {
		go busi(ch, i)
	}
	//业务放入channel
	for i := 0; i < task_cnt; i++ {
		startTask(ch, i)
	}

	wg.Wait()
}
