package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/douglas-soares/rpc-quick/comparison/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct{}

func (s *Server) Fibonacci(ctx context.Context, req *proto.FiboRequest) (*proto.FiboResponse, error) {
	return &proto.FiboResponse{
		Result: fibonacci(req.Value),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Failed to start listening: %v", err)
	}

	s := Server{}

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(GenerateTLSConfig())))
	proto.RegisterFiboServiceServer(grpcServer, &s)

	fmt.Println("listening 8082")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("serve closed: %v", err)
	}
}

func fibonacci(n int64) int64 {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func GenerateTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("../../../cert.pem", "../../../key.pem")
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
