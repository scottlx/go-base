package main

func MyRecover() interface{} {
	return recover()
}

/**

  panic: panic内置函数停止当前goroutine的正常执行，当函数F调用panic时，函数F的正常执行被立即停止，
    然后运行所有在F函数中的defer函数，然后F返回到调用他的函数对于调用者G
  recover: recover内置函数用来管理含有panic行为的goroutine，recover运行在defer函数中，获取panic抛出的错误值，并将程序恢复成正常执行的状态。
    如果在defer函数之外调用recover，那么recover不会停止并且捕获panic错误如果goroutine中没有panic或者捕获的panic的值为nil，recover的返回值也是nil。
    由此可见，recover的返回值表示当前goroutine是否有panic行为

必须要和有异常的栈帧只隔一个栈帧，recover函数才能正常捕获异常。
换言之，recover函数捕获的是祖父一级调用函数栈帧的异常（刚好可以跨越一层defer函数）！


defer 表达式的函数如果定义在 panic 后面，此时panic会在defer还没声明的情况下去调defer导致该函数在 panic 后就无法被执行到


**/
func main() {
	// 嵌套一层函数可以正常捕获异常
	defer MyRecover()
	// 或者用匿名函数
	/*
		defer func(){
			if r = recover(); r != nil {
				fmt.Println("捕获异常: ",r)
			}
		}()
	*/
	panic(1)

}
