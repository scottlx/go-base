package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//粗粒度使用Mutex
var total1 struct {
	sync.Mutex
	value int
}

func worker1(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i <= 100; i++ {
		total1.Lock()
		total1.value += 1
		total1.Unlock()
	}
}

//atomic包原子操作支持
var total2 uint64

//要等待N个线程完成后再进行下一步的同步操作有一个简单的做法，就是使用sync.WaitGroup来等待一组事件,类似创建C缓存大小的channel进行C个协程并发
func worker2(wg *sync.WaitGroup) {
	defer wg.Done()

	var i uint64
	for i = 0; i <= 100; i++ {
		atomic.AddUint64(&total2, i)
	}
}

//单件模式：确保一个类只有一个实例，并提供一个全局访问点。
type singleton struct{}

var (
	instance    *singleton
	initialized uint32
	mu          sync.Mutex
)

func Instance() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}

	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		defer atomic.StoreUint32(&initialized, 1)
		instance = &singleton{}
	}
	return instance
}

//sync.Once实现单例，Once的实现和上述单例模式相同
var (
	instance2 *singleton
	once      sync.Once
)

func Instance2() *singleton {
	once.Do(func() {
		instance2 = &singleton{}
	})
	return instance2
}

/*
//简化的生产者消费者模型
// 后台线程生成最新的配置信息；前台多个工作者线程获取最新的配置信息。所有线程共享配置信息资源
func latestConfig() {
	var config atomic.Value

	config.Store(loadConfig())

	//后台线程加载最新的配置信息
	go func() {
		for {
			time.Sleep(time.Second)
			config.Store(loadConfig())
		}
	}()

	//处理工作请求的线程获取最新的配置信息
	for i := 0; i < 10; i++ {
		go func() {
			for r := range requests() {
				c := config.Load()
				// ...
			}
		}()
	}

}

//根据channel的特性控制并发goroutine的个数
//对于Channel的第K个接收完成操作发生在第K+C个发送操作完成之前
var limit = make(chan int, 3)

func concurentOfThreeRoutine() {
	for _, w := range work {
		go func() {
			limit <- 1
			w()
			<-limit
		}()
	}
	select {}
}
*/

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go worker1(&wg)
	go worker1(&wg)
	wg.Wait()

	fmt.Println(total1.value)
}
