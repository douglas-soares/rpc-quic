package rpc

import "fmt"

type proxy struct {
	serverAddr string
	requestor
}

// newClientProxy creates a client proxy that translates the local invocation
// into parameters for the REQUESTOR, triggers the invocation and returns the result
func newClientProxy(serverAddr string, crh *clientRequestHandler) *proxy {
	r := newRequestor(crh)

	return &proxy{
		serverAddr: serverAddr,
		requestor:  r,
	}
}

func (p *proxy) Call(function string, arg interface{}, result interface{}) error {
	reqResponse := p.Invoke(p.serverAddr, function, arg)
	response := reqResponse.(clientResponse)

	if response.Err != "" {
		return fmt.Errorf(response.Err)
	}

	if result == nil {
		return nil
	}

	return unmarshal(*response.Result, &result)
}
