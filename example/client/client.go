package main

import (
	"crypto/tls"
	"fmt"

	naming "github.com/douglas-soares/rpc-quic/src/naming_service"
	"github.com/douglas-soares/rpc-quic/src/rpc"
	"github.com/lucas-clemente/quic-go"
)

func main() {
	tlsConfN := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"naming"},
	}

	// create the naming service with its known location
	n := naming.NewNamingService("localhost:4040")

	// start naming client
	n.StartClient(tlsConfN)

	// look up for the desired server location
	serverAddr, err := n.LookUp("Servidor")
	if err != nil {
		panic(err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"server"},
		ClientSessionCache: tls.NewLRUClientSessionCache(100), // enables 0-RTT
	}

	//	tokenStore := quic.NewLRUTokenStore(10, 10)
	client := rpc.NewClient(serverAddr, tlsConf, &quic.Config{})

	var resp int

	loop := 1
	for i := 0; i < loop; i++ {

		// call fibonacci function with arg 1
		err := client.Call("fibonacci", 1, &resp)
		if err != nil {
			fmt.Println("client error:", err)
		}

		// client.Close()

		fmt.Println(i, "Client result:", resp)
	}
}
