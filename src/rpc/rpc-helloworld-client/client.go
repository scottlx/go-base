package main

import (
	"fmt"
	helloworldrpc "go-base/src/rpc/helloworld-rpc"
	"log"
	"net/rpc"
)

type HelloServiceClient struct {
	*rpc.Client
}

var _ helloworldrpc.HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}

func (p *HelloServiceClient) Hello(request string, reply *string) error {
	// rpc.Client.Call是同步阻塞调用，首先通过Client.Go方法进行一次异步调用，返回一个表示这次调用的Call结构体。
	// 然后等待Call结构体的Done管道返回调用结果。
	return p.Client.Call(helloworldrpc.HelloServiceName+".Hello", request, reply)
}

// 使用rpc.Client.Go进行异步非阻塞调用，异步调用的输入参数和返回值暂存在rpc.Call中
// rpc.Client.Go通过client.send将call的完整参数发送到RPC框架
// client.send里面有加锁，是线程安全的
// 当调用完成或者发生错误时，将调用call.done方法通知完成, 将Call写入call.Done channel中
func doClientWork(client *rpc.Client) {
	helloCall := client.Go("HelloService.Hello", "hello", new(string), nil)

	// do some thing

	helloCall = <-helloCall.Done
	if err := helloCall.Error; err != nil {
		log.Fatal(err)
	}

	args := helloCall.Args.(string)
	reply := helloCall.Reply.(string)
	fmt.Println(args, reply)
}

func main() {
	client, err := DialHelloService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	err = client.Hello("hello", &reply)
	if err != nil {
		log.Fatal(err)
	}
}
