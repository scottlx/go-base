package main

import (
	"fmt"
	"unsafe"
)

func main() {
	a := [3]int8{6, 8, 9}
	a_first_point := &a[0]
	a_first_unsafe_point := unsafe.Pointer(a_first_point)
	fmt.Println("a[0]的地址为：", a_first_unsafe_point)
	fmt.Println("a[1]的地址为：", unsafe.Pointer(&a[1]))
	a_uintptr_first_unsafe_point := uintptr(a_first_unsafe_point)
	a_uintptr_first_unsafe_point++
	fmt.Printf("a[0]位置指针自增1后的指针位置: %x\n", a_uintptr_first_unsafe_point)
	a_uintptr_second_unsafe_point := unsafe.Pointer(a_uintptr_first_unsafe_point)
	fmt.Println("a[0]位置指针自增1后，的指针位置，转成unsafe_Pointer类型：", a_uintptr_second_unsafe_point)
	int8_point := (*int8)(a_uintptr_second_unsafe_point)
	fmt.Println(*int8_point)
}
