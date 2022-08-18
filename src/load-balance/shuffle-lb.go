package main

import (
	"fmt"
	"math/rand"
	"time"
)

var endpoints = []string{
	"100.69.62.1:3232",
	"100.69.62.32:3232",
	"100.69.62.42:3232",
	"100.69.62.81:3232",
	"100.69.62.11:3232",
	"100.69.62.113:3232",
	"100.69.62.101:3232",
}

// fisher-yates算法
func shuffle(n int) []int {
	rand.Seed(time.Now().UnixNano())
	b := rand.Perm(n)
	return b
}

func request(params map[string]interface{}) error {
	var err error
	indexes := shuffle(len(endpoints))
	maxRetryTimes := 3

	idx := 0
	for i := 0; i < maxRetryTimes; i++ {
		err = apiRequest(params, indexes[idx])
		if err == nil {
			break
		}
		fmt.Println(err)
		idx++
	}

	if err != nil {
		// logging
		fmt.Printf("Failed with %d tries\n", maxRetryTimes)
		return err
	}

	return nil
}
func apiRequest(params map[string]interface{}, index int) error {
	f := params[endpoints[index]].(func(name string) error)
	err := f(endpoints[index])
	return err
}

func main() {
	okFunc := func(name string) error {
		fmt.Println(name)
		return nil
	}
	errFunc := func(name string) error {
		return fmt.Errorf("server:%v not available", name)
	}
	params := make(map[string]interface{})
	for i, endpoint := range endpoints {
		if i%2 == 0 {
			params[endpoint] = errFunc
		} else {
			params[endpoint] = okFunc
		}
	}
	request(params)

	// 均衡验证
	var cnt = map[int]int{}
	for i := 0; i < 1000000; i++ {
		var sl = []int{0, 1, 2, 3, 4, 5, 6}
		sl = shuffle(len(sl))
		cnt[sl[0]]++
	}
	fmt.Println(cnt)
}
