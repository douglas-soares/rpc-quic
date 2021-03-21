package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	naming "github.com/douglas-soares/rpc-quick/src/naming_service"
	"github.com/douglas-soares/rpc-quick/src/rpc"
)

func main() {
	gob.Register(naming.NamingResult{})

	tlsConfN := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"naming"},
	}
	n := naming.NewNamingService("localhost:4040")
	n.StartClient(tlsConfN)
	s, err := n.LookUp("Servidor")
	fmt.Println(" testing naming:", s, err)
	cert, err := tls.LoadX509KeyPair("../../cert.pem", "../../key.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCert, err := ioutil.ReadFile("../../cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConf := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-echo-example"},
	}

	proxy := rpc.NewClient(s.Addr, tlsConf, nil)
	start := time.Now()
	for i := 0; i < 5000; i++ {
		var resp int
		err = proxy.Call(&resp, "fibonacci", 20)
		if err != nil {
			fmt.Println("client error:", err)
		}

		fmt.Println(i, "Client result:", resp)
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed.Milliseconds())
}
