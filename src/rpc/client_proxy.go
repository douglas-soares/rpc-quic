package rpc

import (
	"fmt"
	"reflect"
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

	if response.Err != nil {
		return response.Err
	}

	if result == nil {
		return nil
	}

	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer of interface to receive the rpc result value")
	}

	resultType := reflect.Indirect(reflect.ValueOf(result)).Type()
	responseType := reflect.TypeOf(response.Result)

	if resultType != responseType {
		return fmt.Errorf(fmt.Sprintf("different types: type of result is %s and the response is %s", resultType.Name(), responseType.Name()))
	}

	value := reflect.ValueOf(response.Result)
	reflect.Indirect(reflect.ValueOf(result)).Set(value)

	return nil

}
