package helloworldrpc

import "net/rpc"

// 接口规范分为三个部分：首先是服务的名字，
const HelloServiceName = "path/to/pkg.HelloService"

//然后是服务要实现的详细方法列表
type HelloServiceInterface = interface {
	Hello(request string, reply *string) error
}

//最后是注册该类型服务的函数
func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}

//标准库的RPC默认采用Go语言特有的gob编码，不能跨语言调用
