package main

import (
	"context"
	"log"
	"time"

	pb "github.com/sarthak0714/dbz/pkg/api"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDatabaseClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Simple Put operation
	_, err = c.Put(ctx, &pb.PutRequest{Key: "hello", Value: "world"})
	if err != nil {
		log.Fatalf("could not put: %v", err)
	}
	log.Println("Put successful")

	// Simple Get operation
	getResp, err := c.Get(ctx, &pb.GetRequest{Key: "hello"})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}
	log.Printf("Get result: %s", getResp.Value)

	// Transaction example
	prepareResp, err := c.PrepareTransaction(ctx, &pb.PrepareRequest{Keys: []string{"key1", "key2"}})
	if err != nil {
		log.Fatalf("could not prepare transaction: %v", err)
	}
	if !prepareResp.Ready {
		log.Fatalf("transaction not ready")
	}
	log.Println("Transaction prepared")

	commitResp, err := c.CommitTransaction(ctx, &pb.CommitRequest{
		Updates: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})
	if err != nil {
		log.Fatalf("could not commit transaction: %v", err)
	}
	if !commitResp.Success {
		log.Fatalf("transaction commit failed")
	}
	log.Println("Transaction committed")

	// Verify transaction results
	getResp1, err := c.Get(ctx, &pb.GetRequest{Key: "key1"})
	if err != nil {
		log.Fatalf("could not get key1: %v", err)
	}
	log.Printf("Get key1 result: %s", getResp1.Value)

	getResp2, err := c.Get(ctx, &pb.GetRequest{Key: "key2"})
	if err != nil {
		log.Fatalf("could not get key2: %v", err)
	}
	log.Printf("Get key2 result: %s", getResp2.Value)

	// Delete operation
	_, err = c.Delete(ctx, &pb.DeleteRequest{Key: "hello"})
	if err != nil {
		log.Fatalf("could not delete: %v", err)
	}
	log.Println("Delete successful")

	// Verify delete
	_, err = c.Get(ctx, &pb.GetRequest{Key: "hello"})
	if err == nil {
		log.Fatalf("key 'hello' should have been deleted")
	}
	log.Println("Verified delete: key 'hello' not found as expected")
}
