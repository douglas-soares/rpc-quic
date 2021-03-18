package main

import (
	"fmt"
	"net/rpc"
	"time"
)

func main() {
	client, _ := rpc.Dial("tcp", "localhost:6566")
	start := time.Now()
	for i := 0; i < 1000; i++ {
		var resp int
		err := client.Call("Servidor.Fibonacci", 20, &resp)
		if err != nil {
			fmt.Println("client error:", err)
		}

		fmt.Println(i, "Client result:", resp)

	}
	elapsed := time.Since(start)
	fmt.Println(elapsed.Milliseconds())
}
