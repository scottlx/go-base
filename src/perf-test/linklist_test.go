package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"golang.org/x/sync/errgroup"
)

type Node struct {
	Value interface{}
	Next  *Node
}

// WithLockList
type WithLockList struct {
	Head *Node
	mu   sync.Mutex
}

func (l *WithLockList) Push(v interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := &Node{
		Value: v,
		Next:  l.Head,
	}
	l.Head = n
}

func (l WithLockList) String() string {
	s := ""
	cur := l.Head
	for {
		if cur == nil {
			break
		}
		if s != "" {
			s += ","
		}
		s += fmt.Sprintf("%v", cur.Value)
		cur = cur.Next
	}
	return s
}

type LockFreeList struct {
	Head atomic.Value
}

func (l *LockFreeList) Push(v interface{}) {
	for {
		//保存当前head
		head := l.Head.Load()
		headNode, _ := head.(*Node)
		n := &Node{
			Value: v,
			Next:  headNode,
		}
		//若head在更新时被其他goroutine改变，则更新失败
		if l.Head.CompareAndSwap(head, n) {
			break
		}
	}
}

func (l LockFreeList) String() string {
	s := ""
	cur := l.Head.Load().(*Node)
	for {
		if cur == nil {
			break
		}
		if s != "" {
			s += ","
		}
		s += fmt.Sprintf("%v", cur.Value)
		cur = cur.Next
	}
	return s
}

// ConcurWriteWithLockList
func ConcurWriteWithLockList(l *WithLockList) {
	var g errgroup.Group

	for i := 0; i < 10; i++ {
		i := i
		g.Go(func() error {
			l.Push(i)
			return nil
		})
	}
	_ = g.Wait()
}

// ConcurWriteLockFreeList
func ConcurWriteLockFreeList(l *LockFreeList) {
	var g errgroup.Group
	for i := 0; i < 10; i++ {
		i := i
		g.Go(func() error {
			l.Push(i)
			return nil
		})
	}
	_ = g.Wait()
}

func BenchmarkWriteWithLockList(b *testing.B) {
	l := &WithLockList{}
	for n := 0; n < b.N; n++ {
		l.Push(n)
	}
}

func BenchmarkWriteLockFreeList(b *testing.B) {
	l := &LockFreeList{}
	for n := 0; n < b.N; n++ {
		l.Push(n)
	}
}

func main() {
	l1 := &WithLockList{}
	ConcurWriteWithLockList(l1)
	fmt.Println(l1)
	l2 := &LockFreeList{}
	ConcurWriteLockFreeList(l2)
	fmt.Println(l2)
}
