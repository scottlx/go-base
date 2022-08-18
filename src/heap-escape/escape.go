package main

import "fmt"

type A struct {
	s string
}

// "在方法内把局部变量指针返回" 的情况
func foo(s string) *A {
	a := new(A)
	a.s = s
	return a //返回局部变量a,在C语言中妥妥野指针，但在go则ok，但a会逃逸到堆
}

// go build -gcflags=-m
func main() {
	a := foo("hello")
	b := a.s + " world"
	c := b + "!"
	//interface{}类型一般情况下底层会进行reflect，而使用的reflect.TypeOf(arg).Kind()获取接口类型对象的底层数据类型时发生了堆逃逸
	fmt.Println(c)
}
