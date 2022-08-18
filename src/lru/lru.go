package main

import (
	"container/list"
	"fmt"
)

type Cache struct {
	MaxEntries int
	ll         *list.List
	cache      map[interface{}]*list.Element
	OnEvicted  func(key Key, value interface{})
}

type Key interface{}

type entry struct {
	key   Key
	value interface{}
}

func New(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

func (c *Cache) Add(key Key, value interface{}) {
	// 未创建初始化
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	// 若key已经cache，则刷到队头
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		e.Value.(*entry).value = value
		return
	}
	// 若首次添加，加入map，并刷到队头
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele

	// 若添加后超过缓存大小，删掉最老的
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	// 每次访问后, 刷新被访问的元素到队头
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

func (c *Cache) Remove(key Key) {
	if c.cache == nil {
		return
	}

	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

func (c *Cache) Clear() {
	for _, e := range c.cache {
		kv := e.Value.(*entry)
		c.OnEvicted(kv.key, kv.value)
	}
	c.cache = nil
	c.ll = nil
}

func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

func main() {
	c := New(5)
	c.Add("a", 1)
	c.Add("b", 2)
	c.Add("c", 3)
	c.Add("d", 4)
	c.Add("e", 5)
	c.Get("a")
	c.Add("f", 6)
	c.Get("a")
	c.RemoveOldest()
	if v, ok := c.Get("a"); ok {
		fmt.Printf("got a : %v\n", v)
	} else {
		fmt.Println("a does not exist")
	}
}
