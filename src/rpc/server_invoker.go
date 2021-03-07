package rpc

import (
	"fmt"
	"reflect"
)

type invoker struct {
	funcs map[string]reflect.Value
}

func newInvoker() *invoker {
	funcs := make(map[string]reflect.Value)
	return &invoker{funcs: funcs}
}

// Register the name of the function and its entries
func (i *invoker) Register(function string, fFunc interface{}) {
	if _, ok := i.funcs[function]; ok {
		return
	}
	i.funcs[function] = reflect.ValueOf(fFunc)
	fmt.Println("fucntion", function, "registred")
}

// fazer funcao de unregister

func (i *invoker) invoke(data []byte) []byte {
	req, err := unmarshall(data)
	if err != nil {
		fmt.Println("invoker 1", err)
		return i.returnError(err)
	}

	response := i.execute(req)

	marshalledResponse, err := marshall(response)
	if err != nil {
		fmt.Println("invoker 2", err)
		return i.returnError(err)
	}

	return marshalledResponse
}

func (i *invoker) execute(req rpcData) rpcData {
	f, ok := i.funcs[req.Function]
	if !ok {
		err := fmt.Errorf("func %s not registered", req.Function)
		return rpcData{Args: nil, Err: err}
	}

	//log.Printf("func %s is called\n", req.Function)
	// unpackage request arguments
	inArgs := make([]reflect.Value, len(req.Args))
	for i := range req.Args {
		inArgs[i] = reflect.ValueOf(req.Args[i])
	}
	// invoke requested method
	out := f.Call(inArgs)
	// now since we have followed the function signature style where last argument will be an error
	// so we will pack the response arguments expect error.

	if len(out) == 0 {
		return rpcData{}
	}
	resArgs := make([]interface{}, len(out))
	for i := 0; i < len(out); i++ {
		// Interface returns the constant value stored in v as an interface{}.
		resArgs[i] = out[i].Interface()
	}

	// // pack error argument
	// var err error
	// if e, ok := resArgs[len(out)-1].(error); ok {
	// 	// convert the error into error string value
	// 	resArgs = resArgs[:len(out)-1]
	// 	err = e
	// }
	return rpcData{Result: resArgs[0]} // fix this later
}

func (i *invoker) returnError(err error) []byte {
	resp := rpcData{
		Err: err,
	}
	r, _ := marshall(resp) // retornar erro para cancelar conexao no servidor
	return r
}
