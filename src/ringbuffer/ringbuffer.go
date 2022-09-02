package ringbuffer

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Ring struct {
	consHead uint32
	consTail uint32
	prodHead uint32
	prodTail uint32
	mem      []string
	length   uint32
}

func NewRing(size uint32) *Ring {
	return &Ring{
		consHead: 0,
		consTail: 0,
		prodHead: 0,
		prodTail: 0,
		mem:      make([]string, size),
		length:   size,
	}
}

func (r *Ring) Enqueue(elem string) {
	err := r.enqueue(elem)
	for err != nil {
		err = r.enqueue(elem)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//overwrite mode
func (r *Ring) enqueue(elem string) error {
	var localProdNext, localProdHead, localConsTail uint32
	localProdHead = atomic.LoadUint32(&r.prodHead)
	localConsTail = atomic.LoadUint32(&r.consTail)
	if localProdHead < r.length-1 {
		localProdNext = localProdHead + 1
	}

	if localConsTail > localProdHead {
		return errors.New("buffer is full")
	}

	if success := atomic.CompareAndSwapUint32(&r.prodHead, localProdHead, localProdNext); success {
		r.mem[localProdHead] = elem
		fmt.Println(elem, " writted")
		for r.prodTail != localProdHead {
			//fmt.Println("waiting others finish writing")
		}
		r.prodTail = localProdNext
		fmt.Println(r)
		return nil
	} else {
		fmt.Println("preempted by others")
		return errors.New("preempted by others")
	}
}

func (r *Ring) Dequeue() (string, error) {
	var localConsHead, localConsNext, localProdTail uint32
	var res string

	localConsHead = atomic.LoadUint32(&r.consHead)
	localProdTail = atomic.LoadUint32(&r.prodTail)
	if localConsHead != r.length-1 {
		localConsNext = localConsHead + 1
	}

	if localConsHead == localProdTail {
		return "", errors.New("buffer is empty")
	}

	r.consHead = localConsNext
	res = r.mem[localConsHead]
	r.consTail = r.consHead

	return res, nil
}

func TestSingleRoutine(t *testing.T) {
	r := NewRing(5)
	var wg sync.WaitGroup
	wg.Add(10)

	var locker = new(sync.Mutex)
	var cond = sync.NewCond(locker)

	for i := 0; i < 5; i++ {
		go func(elem int) {
			defer wg.Done()

			cond.L.Lock()
			cond.Wait()

			r.Enqueue(strconv.Itoa(elem))

			cond.L.Unlock()

		}(i)
	}

	time.Sleep(time.Second * 1)
	cond.Broadcast()
	wg.Wait()
}

func TestMultiRoutine(t *testing.T) {
	r := NewRing(5)
	for i := 0; i < 6; i++ {
		r.Enqueue(strconv.Itoa(i))
	}
	for i := 0; i < 5; i++ {
		if str, err := r.Dequeue(); err == nil {
			fmt.Println("get ", str)
			fmt.Println(r)
		} else {
			fmt.Println("error: ", err)
			fmt.Println(r)
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
