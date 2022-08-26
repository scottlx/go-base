package main

import (
	"fmt"
	"go-base/src/ringbuffer"
	"strconv"
)

func main() {
	r := ringbuffer.NewRing(5)
	for i := 0; i < 4; i++ {
		r.Enqueue(strconv.Itoa(i))
	}
	for i := 0; i < 5; i++ {
		if str, err := r.Dequeue(); err == nil {
			fmt.Println("get ", str)
		} else {
			fmt.Println("error: ", err)
		}
	}
	r.Enqueue(strconv.Itoa(5))
	if str, err := r.Dequeue(); err == nil {
		fmt.Println("get ", str)
	} else {
		fmt.Println("error: ", err)
	}

	//ringbuffer.Test()
}
