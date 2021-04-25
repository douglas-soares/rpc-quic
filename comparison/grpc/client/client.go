package main

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/douglas-soares/rpc-quick/comparison/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := grpc.Dial(":8082", grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
	if err != nil {
		panic(err)
	}

	loop := 1000
	for i := 0; i < loop; i++ {
		client := proto.NewFiboServiceClient(conn)

		response, err := client.Fibonacci(context.Background(), &proto.FiboRequest{Value: 1})
		if err != nil {
			panic(err)
		}

		fmt.Println(i, "client result:", response.Result)
	}
}
