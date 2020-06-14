package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"kv-ttl/kv"
	"kv-ttl/pb"
	"kv-ttl/repository"
	"kv-ttl/repository/postgres"
	"kv-ttl/server"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	cacheConfig := kv.Configuration{
		BackupInterval: backupInterval(),
		Storage:        storage(),
	}
	cache := kv.NewCache(cacheConfig)
	cacheServer := server.NewCacheServer(cache)

	opts := make([]grpc.ServerOption, 0)
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterStorageServer(grpcServer, cacheServer)

	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(listener)
}

func backupInterval() time.Duration {
	biEnv := os.Getenv("BP_INTERVAL")
	if biEnv != "" {
		return kv.DefaultBackupInterval
	}
	bi, err := strconv.Atoi(biEnv)
	if err != nil {
		return kv.DefaultBackupInterval
	}
	return time.Duration(bi)
}

// storage parses environment variables and configures one of supported data storages.
func storage() kv.Storage {
	switch os.Getenv("STORAGE") {
	case "pg":
		port, err := strconv.Atoi(os.Getenv("PG_PORT"))
		if err != nil {
			panic("cannot parse db port: " + err.Error())
		}
		connection := postgres.Connection{
			Host:     os.Getenv("PG_HOST"),
			Port:     port,
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PWD"),
			Database: os.Getenv("PG_DB"),
		}
		db, err := postgres.NewPostgresDb(connection)
		if err != nil {
			panic(err)
		}
		fmt.Printf("started with database storage => %s:%d\n", connection.Host, connection.Port)
		return postgres.NewRepository(db)

	case "file":
		filename := os.Getenv("FNAME")
		fmt.Printf("started with file storage => %s\n", filename)
		return repository.NewFileRepo(filename)

	default:
		fmt.Println("started without persistent storage")
		return &kv.UnimplementedStorage{}
	}
}
