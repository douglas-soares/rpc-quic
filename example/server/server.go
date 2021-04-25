package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"

	naming "github.com/douglas-soares/rpc-quick/src/naming_service"
	"github.com/douglas-soares/rpc-quick/src/rpc"
	quic "github.com/lucas-clemente/quic-go"
)

func main() {
	tlsConfN := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"naming"},
	}

	// create the naming service with its known location
	n := naming.NewNamingService("localhost:4040")

	// start the naming client
	n.StartClient(tlsConfN)

	// register the server location
	serverAddr := "localhost:8080"
	err := n.Bind("Servidor", serverAddr)
	if err != nil {
		panic(err)
	}

	server := rpc.NewServer()
	// register server remote function
	server.Register("fibonacci", fibonacci)

	// start to listen the port
	server.ListenAndServe(serverAddr, generateTLSConfig(), &quic.Config{})
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// Setup a bare-bones TLS config for the server
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
		NextProtos:   []string{"server"},
	}
}
