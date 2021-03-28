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
		InsecureSkipVerify: true,
		//RootCAs:            caCertPool,
		NextProtos:         []string{"quic-echo-example"},
		ClientSessionCache: tls.NewLRUClientSessionCache(100),
	}

	//	tokenStore := quic.NewLRUTokenStore(10, 10)
	quicConfig := &quic.Config{}

	client := rpc.NewClient("localhost:8080", tlsConf, quicConfig)
	start := time.Now()
	total := float64(0)
	loop := 1
	for i := 0; i < loop; i++ {
		t0 := time.Now()

		var resp int
		err := client.Call("fibonacci", 1, &resp)
		if err != nil {
			fmt.Println("client error:", err)
		}

		client.Close()
		t1 := time.Since(t0)
		total = total + float64(t1.Milliseconds())
		fmt.Println(i, "Client result:", resp)
	}
	elapsed := time.Since(start)
	fmt.Println("Total:", elapsed.Milliseconds())
	fmt.Println("Mean", total/float64(loop))

}
