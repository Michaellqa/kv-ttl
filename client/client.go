// The client demonstrates how to connect to grpc server and call the cache methods
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"kv-ttl/kv"
	"kv-ttl/pb"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:80", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cl := pb.NewStorageClient(conn)

	ctx := context.Background()
	_, err = cl.Add(ctx, &pb.KeyValue{Key: "1", Value: &pb.T{Value: "One"}})
	_, err = cl.Add(ctx, &pb.KeyValue{Key: "2", Value: &pb.T{Value: "Two"}})
	_, err = cl.Add(ctx, &pb.KeyValue{Key: "3", Value: &pb.T{Value: "Three"}})

	resp, err := cl.Get(ctx, &pb.Key{Key: "0"})
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("#0: %v\n", resp.Value)
	}

	resp, err = cl.Get(ctx, &pb.Key{Key: "1"})
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("#1: %v\n", resp.Value)
	}

	for {
		stream, err := cl.GetAll(ctx, &pb.Empty{})
		if err != nil {
			log.Println(err)
			return
		}
		var values []kv.T
		for {
			value, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			values = append(values, kv.T{V: value.Value})
		}
		fmt.Printf("#all: %v", values)
		time.Sleep(2 * time.Second)
	}
}
