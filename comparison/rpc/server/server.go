package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type server struct{}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
func (s server) Fibonacci(n int, result *int) error {
	f := fibonacci(n)
	*result = f
	return nil
}

func Server() {

	server2 := rpc.NewServer()
	server2.RegisterName("Servidor", server{})
	l, _ := net.Listen("tcp", "localhost:6566")
	fmt.Println("Iniciando conexão...")
	server2.Accept(l)
}

func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go Server()
	waitGroup.Wait()
}
