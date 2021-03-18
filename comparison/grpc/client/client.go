package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/douglas-soares/rpc-quick/comparison/grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := proto.NewFiboServiceClient(conn)

	start := time.Now()
	for i := 0; i < 10000; i++ {
		response, err := client.Fibonacci(context.Background(), &proto.FiboRequest{Value: 20})
		if err != nil {
			log.Fatalf("Error calling fibonacci method: %v", err)
		}
		fmt.Println(i, "client result:", response.Result)
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed.Milliseconds())
}
