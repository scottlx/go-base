package main

import (
	"fmt"
	"time"
)

const (
	fmat = "2006-01-02 15:04:05"
)

func ReadFromClosedChan() {
	c := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		c <- 10
		close(c)
	}()
	for {
		select {
		case x, ok := <-c:
			fmt.Printf("%v, read channel: x=%v, ok=%v\n", time.Now().Format(fmat), x, ok)
			time.Sleep(500 * time.Millisecond)

		default:
			fmt.Printf("%v, go into Default\n", time.Now().Format(fmat))
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func DontReadFromCloseChan() {
	c := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		c <- 10
		close(c)
	}()
	for {
		select {
		case x, ok := <-c:
			fmt.Printf("%v, read channel: x=%v, ok=%v\n", time.Now().Format(fmat), x, ok)
			time.Sleep(500 * time.Millisecond)
			if !ok {
				c = nil
			}

		default:
			fmt.Printf("%v, go into Default\n", time.Now().Format(fmat))
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func CloseChanWithoutDefault() {
	c := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		c <- 10
		close(c)
	}()
	for {
		select {
		case x, ok := <-c:
			fmt.Printf("%v, read channel: x=%v, ok=%v\n", time.Now().Format(fmat), x, ok)
			time.Sleep(500 * time.Millisecond)
			// change closed chan to nil without default statement will cause dead lock !!
			/*
				if !ok {
					c = nil
				}
			*/
		}
	}
}

func main() {
	CloseChanWithoutDefault()
}
