package main

import (
	"fmt"
	"google.golang.org/grpc"
	"kv-ttl/kv"
	"kv-ttl/pb"
	"kv-ttl/server"
	"log"
	"net"
)

type Foo func() int

func main() {
	port := 80
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	cache := kv.NewCache(kv.Configuration{})
	cacheServer := server.NewCacheServer(cache)

	pb.RegisterStorageServer(grpcServer, cacheServer)
	grpcServer.Serve(listener)
}
