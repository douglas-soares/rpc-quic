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

	proxy := rpc.NewClient(s.Addr, tlsConf)
	start := time.Now()
	for i := 0; i < 1; i++ {
		var resp int
		err = proxy.Call(&resp, "fibonacci", 10)
		if err != nil {
			fmt.Println("client error:", err)
		}

		fmt.Println(i, "Client result:", resp)

	}
	elapsed := time.Since(start)
	fmt.Println(elapsed.Milliseconds())
}
