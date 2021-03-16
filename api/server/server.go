package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	naming "github.com/douglas-soares/rpc-quick/src/naming_service"
	"github.com/douglas-soares/rpc-quick/src/rpc"
)

func main() {
	tlsConfN := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"naming"},
	}
	n := naming.NewNamingService("localhost:4040")
	n.StartClient(tlsConfN)
	err := n.Bind("Servidor", "localhost:4242")
	fmt.Println("server testing naming:", err)

	server := rpc.NewServer()
	server.Register("fibonacci", fibonacci)

	server.ListenAndServe("localhost:4242", GenerateTLSConfig(), nil)
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func GenerateTLSConfig() *tls.Config {
	// cer, err := tls.LoadX509KeyPair("../../server.crt", "../../server.key")
	// if err != nil {
	// 	log.Println(err)
	// }
	cert, err := tls.LoadX509KeyPair("../../cert.pem", "../../key.pem")
	caCert, err := ioutil.ReadFile("../../cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tls := &tls.Config{
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-echo-example"},
	}

	tls.BuildNameToCertificate()
	return tls
}
