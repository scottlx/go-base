package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	//替代rpc.Client，传入codec
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var reply string

	// {"method":"HelloService.Hello","params":["hello"],"id":0}
	// id是由调用端维护的一个唯一的调用编号
	err = client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
