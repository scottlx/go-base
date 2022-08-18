package main

import (
	"log"
	"sync"
	"time"
)

var done = false

func read(name string, c *sync.Cond) {
	c.L.Lock()
	for !done {
		//调用wait时会自动释放锁，并挂起调用者所在的goroutine
		//收到broadcast之后会重新加锁
		c.Wait()
	}
	log.Println(name, "starts reading")
	c.L.Unlock()
}
func write(name string, c *sync.Cond) {
	log.Println(name, "starts writing")
	//确保三个协程都处于wait状态
	time.Sleep(time.Second)
	done = true
	log.Println(name, "wakes all")
	c.Broadcast()
}
func main() {
	cond := sync.NewCond(&sync.Mutex{})
	go read("reader1", cond)
	go read("reader2", cond)
	go read("reader3", cond)
	write("writer", cond)
	time.Sleep(time.Second * 3)
}
