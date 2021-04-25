package rpc

import "fmt"

type proxy struct {
	serverAddr string
	requestor
}

func newClientProxy(serverAddr string, crh *clientRequestHandler) *proxy {
	r := newRequestor(crh)

	return &proxy{
		serverAddr: serverAddr,
		requestor:  r,
	}
}

func (p *proxy) Call(function string, arg interface{}, result interface{}) error {
	response := p.invoke(p.serverAddr, function, arg)

	if response.Err != "" {
		return fmt.Errorf(response.Err)
	}

	if result == nil {
		return nil
	}

	return unmarshal(*response.Result, &result)
}
