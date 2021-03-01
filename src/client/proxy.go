package client

import "github.com/douglas-soares/rpc-quick/src/models"

// Proxy contains the method for the client proxy
type Proxy interface {
	Call(result interface{}, function string, args []interface{}) error
}

type proxy struct {
	Requestor
	lookupLocation string
}

// NewClientProxy creates a client proxy that translates the local invocation
// into parameters for the REQUESTOR, triggers the invocation and returns the result
func NewClientProxy(lookupLocation string) Proxy {
	var h requestHandler // temporario
	r := NewRequestor(h)

	return &proxy{
		Requestor:      r,
		lookupLocation: lookupLocation,
	}
}

func (p *proxy) Call(result interface{}, function string, args []interface{}) error {
	// lookup
	var location string // resultado do lookup
	reqResponse := p.Invoke(location, function, args)

	response := reqResponse.(models.Response)

	// Usar decode para transformar o Response.Result na mesma interface do result

	return response.Err
}
