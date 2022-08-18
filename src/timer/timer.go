package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool)

	go func() {

		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 7)
			c <- false
		}

		time.Sleep(time.Second * 7)
		c <- true
	}()

	go func() {
		timer := time.NewTimer(time.Second * 5)
		for {
			// 如果明确time已经expired，并且t.C已经被取空，那么可以直接使用Reset；
			// 如果程序之前没有从t.C中读取过值，这时需要首先调用Stop()，
			// 如果返回true，说明timer还没有expire，stop成功删除timer，可直接reset；
			// 如果返回false，说明stop前已经expire，需要显式drain channel。
			if !timer.Stop() {
				// 无论timer.C种是否有数据，都不会被阻塞
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(time.Second * 5)
			select {
			case b := <-c:
				if b == false {
					fmt.Println(time.Now(), ":recv false. continue")
					continue
				}
				//we want true, not false
				fmt.Println(time.Now(), ":recv true. return")
				return
			case <-timer.C:
				fmt.Println(time.Now(), ":timer expired")
				continue
			}
		}
	}()

	//to avoid that all goroutine blocks.
	var s string
	fmt.Scanln(&s)
}
