package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"kv-ttl/kv"
	"kv-ttl/pb"
	"kv-ttl/repository/postgres"
	"kv-ttl/server"
	"log"
	"net"
	"time"
)

type Foo func() int

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	port := 80
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	//host := "db"
	host := "localhost"
	connection := postgres.Connection{
		Host:     host,
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "cache_app",
	}
	db, err := postgres.NewPostgresDb(connection)
	if err != nil {
		panic(err)
	}
	storage := postgres.NewRepository(db)
	cacheConfig := kv.Configuration{
		BackupInterval: 5 * time.Second,
		Storage:        storage,
	}
	cache := kv.NewCache(cacheConfig)

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	cacheServer := server.NewCacheServer(cache)

	pb.RegisterStorageServer(grpcServer, cacheServer)
	grpcServer.Serve(listener)
}
