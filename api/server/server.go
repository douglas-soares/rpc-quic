package main

import (
	"crypto/tls"

	"github.com/douglas-soares/rpc-quick/src/rpc"
	quic "github.com/lucas-clemente/quic-go"
)

func main() {
	// tlsConfN := &tls.Config{
	// 	InsecureSkipVerify: true,
	// 	NextProtos:         []string{"naming"},
	// }
	// n := naming.NewNamingService("localhost:4040")
	// n.StartClient(tlsConfN)
	// err := n.Bind("Servidor", "localhost:4242")
	// fmt.Println("server testing naming:", err)

	server := rpc.NewServer()
	server.Register("fibonacci", fibonacci)

	quicConfig := &quic.Config{}
	server.ListenAndServe("localhost:8080", GenerateTLSConfig(), quicConfig)
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func GenerateTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("../../cert.pem", "../../key.pem")
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
