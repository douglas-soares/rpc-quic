package rpc

import (
	"fmt"
)

type proxy struct {
	serverAddr string
	requestor
}

// newClientProxy creates a client proxy that translates the local invocation
// into parameters for the REQUESTOR, triggers the invocation and returns the result
func newClientProxy(serverAddr string, crh clientRequestHandler) *proxy {
	r := newRequestor(crh)

	return &proxy{
		serverAddr: serverAddr,
		requestor:  r,
	}
}

func (p *proxy) Call(result interface{}, function string, args ...interface{}) error {
	reqResponse := p.Invoke(p.serverAddr, function, args)
	response := reqResponse.(rpcData)

	fmt.Println("Client proxy response:", response)

	// Usar decode para transformar o Response.Result na mesma interface do result

	return response.Err
}
