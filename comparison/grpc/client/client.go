package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/douglas-soares/rpc-quick/comparison/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	start := time.Now()
	total := float64(0)
	loop := 1000
	for i := 0; i < loop; i++ {
		t0 := time.Now()
		tlsConf := &tls.Config{
			InsecureSkipVerify: true}
		conn, err := grpc.Dial(":8082", grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}

		client := proto.NewFiboServiceClient(conn)

		response, err := client.Fibonacci(context.Background(), &proto.FiboRequest{Value: 1})
		if err != nil {
			log.Fatalf("Error calling fibonacci method: %v", err)
		}
		fmt.Println(i, "client result:", response.Result)
		t1 := time.Since(t0)
		total = total + float64(t1.Milliseconds())
		conn.Close()
	}
	elapsed := time.Since(start)
	fmt.Println("Total:", elapsed.Milliseconds())
	fmt.Println("Mean", total/float64(loop))
}
