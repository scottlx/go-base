package main

import (
	"sync"
	"testing"
	"time"
)

// rpct为读的占比
func OpMapWithMutex(rpct int) {
	m := make(map[int]struct{})
	mu := sync.Mutex{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()

			// 写操作
			if i >= rpct {
				m[i] = struct{}{}
				time.Sleep(time.Microsecond)
				return
			}
			// 读操作
			_ = m[i]
			time.Sleep(time.Microsecond)
		}()
	}
	wg.Wait()
}

func OpMapWithRWMutex(rpct int) {
	m := make(map[int]struct{})
	mu := sync.RWMutex{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 写操作
			if i >= rpct {
				mu.Lock()
				defer mu.Unlock()
				m[i] = struct{}{}
				time.Sleep(time.Microsecond)
				return
			}
			// 读操作
			mu.RLock()
			defer mu.RUnlock()
			_ = m[i]
			time.Sleep(time.Microsecond)
		}()
	}
	wg.Wait()
}

func BenchmarkMutexReadMore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OpMapWithMutex(80)
	}
}
func BenchmarkRWMutexReadMore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OpMapWithRWMutex(80)
	}
}
func BenchmarkMutexRWEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OpMapWithMutex(50)
	}
}
func BenchmarkRWMutexRWEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OpMapWithRWMutex(50)
	}
}
func BenchmarkMutexWriteMore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OpMapWithMutex(20)
	}
}
func BenchmarkRWMutexWriteMore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OpMapWithRWMutex(20)
	}
}
