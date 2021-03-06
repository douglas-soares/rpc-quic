package rpc

import (
	"fmt"
	"log"
	"reflect"
)

type Invoker interface {
	Register(fnName string, fFunc interface{})
}

type invoker struct {
	funcs map[string]reflect.Value
}

func newInvoker() *invoker {
	funcs := make(map[string]reflect.Value)
	return &invoker{funcs: funcs}
}

// Register the name of the function and its entries
func (i *invoker) Register(fnName string, fFunc interface{}) {

	if _, ok := i.funcs[fnName]; ok {
		return
	}
	// colocar thread aqui
	i.funcs[fnName] = reflect.ValueOf(fFunc)
	fmt.Println("fucntion", fnName, "registred")
}

// fazer funcao de unregister

func (i *invoker) invoke(data []byte) []byte {
	fmt.Println(" invoker")
	req, err := unmarshall(data)
	fmt.Println(err)
	if err != nil {
		return i.returnError(err)
	}

	response := i.execute(req)

	marshalledResponse, err := marshall(response)
	if err != nil {
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

	log.Printf("func %s is called\n", req.Function)
	// unpackage request arguments
	inArgs := make([]reflect.Value, len(req.Args))
	for i := range req.Args {
		inArgs[i] = reflect.ValueOf(req.Args[i])
	}
	// invoke requested method
	out := f.Call(inArgs)
	// now since we have followed the function signature style where last argument will be an error
	// so we will pack the response arguments expect error.
	resArgs := make([]interface{}, len(out))
	for i := 0; i < len(out); i++ { // TIREI O -1 DAQUI, TEM QUE TRATAR CORRETAMENTE AQUI
		// Interface returns the constant value stored in v as an interface{}.
		resArgs[i] = out[i].Interface()
	}

	// pack error argument
	var err error
	if e, ok := out[len(out)-1].Interface().(error); ok {
		// convert the error into error string value
		resArgs = resArgs[:len(out)-1] // ta correto isso?
		err = e
	}
	return rpcData{Args: resArgs, Err: err}
}

func (i *invoker) returnError(err error) []byte {
	resp := rpcData{
		Err: err,
	}
	r, _ := marshall(resp) // retornar erro para cancelar conexao no servidor
	return r
}
