package main

import (
	"fmt"
	"net/rpc"
	"time"
)

func main() {

	// tlsConf := &tls.Config{
	// 	InsecureSkipVerify: true,
	// }
	//conn, _ := tls.Dial("tcp", "localhost:6566", tlsConf)
	client, _ := rpc.DialHTTP("tcp", "localhost:6566")
	//client, _ := rpc.Dial("tcp", "localhost:6566")
	start := time.Now()
	total := float64(0)
	loop := 10000

	for i := 0; i < loop; i++ {
		t0 := time.Now()

		var resp int
		err := client.Call("Servidor.Fibonacci", 1, &resp)
		if err != nil {
			fmt.Println("client error:", err)
		}

		fmt.Println(i, "Client result:", resp)

		t1 := time.Since(t0)
		total = total + float64(t1.Milliseconds())
	}
	elapsed := time.Since(start)
	fmt.Println("Total:", elapsed.Milliseconds())
	fmt.Println("Mean", total/float64(loop))
}
