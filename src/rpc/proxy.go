package rpc

import (
	"fmt"
)

// Proxy contains the method for the client proxy
type Proxy interface {
	Call(result interface{}, function string, args ...interface{}) error
}

type proxy struct {
	Requestor
}

// NewClientProxy creates a client proxy that translates the local invocation
// into parameters for the REQUESTOR, triggers the invocation and returns the result
func NewClientProxy() Proxy {
	reqh := newRequestHandler() // deveria ser criado aqui?
	//var h requestHandler        // temporario
	r := NewRequestor(reqh)

	return &proxy{
		Requestor: r,
	}
}

func (p *proxy) Call(result interface{}, function string, args ...interface{}) error {
	// lookup
	//var location string // resultado do lookup
	reqResponse := p.Invoke("", function, args)

	response := reqResponse.(rpcData)

	fmt.Println(response)

	// Usar decode para transformar o Response.Result na mesma interface do result

	return response.Err
}
