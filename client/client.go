// The client demonstrates how to connect to grpc server and call the cache methods
package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"io"
	"kv-ttl/kv"
	"kv-ttl/pb"
	"log"
	"time"
)

const (
	ms  = time.Millisecond
	sec = time.Second
)

// Change this if according to server address
var url = "localhost:80"

// Demonstration of cache features.
func main() {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	cl := pb.NewStorageClient(conn)

	// Populate cache with several values
	_, err = cl.Add(ctx, &pb.KeyValue{Key: "1", Value: pbt("One")})
	_, err = cl.Add(ctx, &pb.KeyValue{Key: "2", Value: pbt("Two")})
	_, err = cl.Add(ctx, &pb.KeyValue{Key: "3", Value: pbt("Three")})
	_, err = cl.AddWithTtl(ctx, &pb.KeyValueTtl{
		Key:   "4",
		Value: pbt("Four"),
		Ttl:   ptypes.DurationProto(10 * time.Second)})
	_, err = cl.AddWithTtl(ctx, &pb.KeyValueTtl{
		Key:   "5",
		Value: pbt("Five"),
		Ttl:   ptypes.DurationProto(3500 * time.Millisecond)})

	// Value existing and nonexistent values
	resp, err := cl.Value(ctx, pbk("5"))
	assertedPrint("Five", resp, err)

	resp, err = cl.Value(ctx, pbk("6"))
	assertedPrint("#6 error not_found", resp, err)

	// Wait until ttl is ended but wasn't swept yet
	time.Sleep(3600 * ms)
	fmt.Println("\nexpected: One Two Three Four Five")
	printAll(cl)

	// Wait until clean completed
	time.Sleep(sec)
	fmt.Println("\nexpected: One Two Three Four")
	printAll(cl)

	// Remove a value
	cl.Remove(ctx, pbk("2"))
	fmt.Println("\nexpected: One Three Four")
	printAll(cl)

	// Value time since value was added
	time.Sleep(2 * time.Second)
	alive, err := cl.TimeAlive(ctx, pbk("1"))
	assertedPrint("#1 alive ~ 6.6 sec", alive, err)

	// Wait until values with ttl are dead.
	// Set ttl to the value that supposed to live forever. Watch him die.
	time.Sleep(6 * sec)
	fmt.Println("\nexpected: One Three")
	printAll(cl)
	stamp, _ := ptypes.TimestampProto(time.Now().Add(1 * sec))
	cl.SetTtl(ctx, &pb.TtlRequest{Key: "3", Stamp: stamp})
	time.Sleep(2 * sec)
	fmt.Println("\nexpected: One")
	printAll(cl)
}

// print conveniently prints to console expected value along with response from the server
func assertedPrint(expected string, v interface{}, err error) {
	fmt.Printf("\nexpected: %s\n", expected)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v)
	}
}

// printAll calls ListAll method and prints all the values
func printAll(cl pb.StorageClient) {
	ctx := context.Background()
	stream, err := cl.ListAll(ctx, &pb.Empty{})
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
	fmt.Printf("#all: %v\n", values)
}

//--- helper functions to initiate protobuf types ---

func pbk(key string) *pb.Key {
	return &pb.Key{Key: key}
}

func pbt(value string) *pb.T {
	return &pb.T{Value: value}
}
