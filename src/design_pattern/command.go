package design_pattern

import (
	"fmt"
	"time"
)

//命令模式

type Cooker struct {
}

func (c *Cooker) MakeChicken() {
	fmt.Println("chicken made")
}

func (c *Cooker) MakeBeef() {
	fmt.Println("Beef made")
}

type Order interface {
	Process()
}

type ChickenOrder struct {
	c Cooker
}

func (o *ChickenOrder) Process() {
	o.c.MakeChicken()
}

type BeefOrder struct {
	c Cooker
}

func (o *BeefOrder) Process() {
	o.c.MakeBeef()
}

type Waiter struct {
	OrderChan chan Order
}

func (w *Waiter) CollectOrder(o Order) {
	w.OrderChan <- o
}

func (w *Waiter) ProcessOrder() {
	for {
		select {
		case o := <-w.OrderChan:
			o.Process()
		}
	}
}

func Run1() {
	cooker := Cooker{}
	o1 := &ChickenOrder{cooker}
	o2 := &BeefOrder{cooker}
	waiter := Waiter{OrderChan: make(chan Order, 5)}
	go func() {
		waiter.ProcessOrder()
	}()

	waiter.CollectOrder(o1)
	waiter.CollectOrder(o2)
	time.Sleep(time.Second)
}
