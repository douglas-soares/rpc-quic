package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
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
		panic(err)
	}

	s := Server{}

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(generateTLSConfig())))
	proto.RegisterFiboServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func fibonacci(n int64) int64 {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}
}
