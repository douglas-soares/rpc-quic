package rpc

import (
	"encoding/json"
)

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

func (p *proxy) Call(result interface{}, function string, args ...interface{}) error {
	reqResponse := p.Invoke(p.serverAddr, function, args)
	response := reqResponse.(rpcData)

	//fmt.Println("Client proxy response:", response)
	if response.Err != nil {
		return response.Err
	}

	b, err := json.Marshal(response.Result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	return nil
}
