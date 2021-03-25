package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/douglas-soares/rpc-quick/src/rpc"
	"github.com/lucas-clemente/quic-go"
)

func main() {
	// gob.Register(naming.NamingResult{})

	// tlsConfN := &tls.Config{
	// 	InsecureSkipVerify: true,
	// 	NextProtos:         []string{"naming"},
	// }
	// n := naming.NewNamingService("localhost:4040")
	// n.StartClient(tlsConfN)
	// s, err := n.LookUp("Servidor")
	// fmt.Println(" testing naming:", s, err)

	caCert, err := ioutil.ReadFile("../../cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConf := &tls.Config{
		RootCAs:            caCertPool,
		NextProtos:         []string{"quic-echo-example"},
		ClientSessionCache: tls.NewLRUClientSessionCache(100),
		ServerName:         "localhost",
	}

	//tokenStore := quic.NewLRUTokenStore(10, 10)
	quicConfig := &quic.Config{}

	client := rpc.NewClient("localhost:8080", tlsConf, quicConfig)
	start := time.Now()
	total := float64(0)
	loop := 5000
	for i := 0; i < loop; i++ {
		t0 := time.Now()

		var resp int
		err := client.Call(&resp, "fibonacci", 1)
		if err != nil {
			fmt.Println("client error:", err)
		}

		t1 := time.Since(t0)
		total = total + float64(t1.Milliseconds())
		fmt.Println(i, "Client result:", resp)
		client.Close()
	}
	elapsed := time.Since(start)
	fmt.Println("Total:", elapsed.Milliseconds())
	fmt.Println("Mean", total/float64(loop))

}
