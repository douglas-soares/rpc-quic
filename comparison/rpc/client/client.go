package main

import (
	"crypto/tls"
	"fmt"
	"net/rpc"
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, _ := tls.Dial("tcp", "localhost:6566", tlsConf)
	client := rpc.NewClient(conn)

	var resp int

	loop := 10000
	for i := 0; i < loop; i++ {
		err := client.Call("Servidor.Fibonacci", 1, &resp)
		if err != nil {
			panic(err)
		}

		fmt.Println(i, "Client result:", resp)
	}
}
