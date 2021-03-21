package main

import (
	"crypto/tls"
	"encoding/gob"
	"fmt"
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

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	client := rpc.NewClient(s.Addr, tlsConf, nil)
	start := time.Now()
	total := float64(0)
	loop := 10000
	for i := 0; i < loop; i++ {
		t0 := time.Now()

		var resp int
		err = client.Call(&resp, "fibonacci", 1)
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
