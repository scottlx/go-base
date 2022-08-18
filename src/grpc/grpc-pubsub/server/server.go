package main

import (
	ps "go-base/src/grpc/grpc-pubsub/pubsub"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	ps.RegisterPubsubServiceServer(grpcServer, new(ps.PubsubService))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
