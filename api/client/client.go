package main

import (
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"time"

	common "github.com/douglas-soares/rpc-quick/api"
	naming "github.com/douglas-soares/rpc-quick/src/naming_service"
	"github.com/douglas-soares/rpc-quick/src/rpc"
)

func main() {
	gob.Register(naming.NamingResult{})
	gob.Register(common.Data{})
	tlsConfN := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"naming"},
	}
	n := naming.NewNamingService("localhost:4040")
	n.StartClient(tlsConfN)
	s, err := n.LookUp("sum")
	fmt.Println(" testing naming:", s, err)

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	proxy := rpc.NewClient(s.Addr, tlsConf)
	start := time.Now()
	for i := 0; i < 5000; i++ {
		var resp common.Data
		err = proxy.Call(&resp, "sum", common.Data{Data: 1.0})
		if err != nil {
			fmt.Println("cliente", err)
		}

		fmt.Println(i, "Client result:", resp.Data)

	}
	elapsed := time.Since(start)
	fmt.Println(elapsed.Milliseconds())
}
