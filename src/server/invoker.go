package server

import (
	"fmt"
	"log"
	"reflect"

	"github.com/douglas-soares/rpc-quick/src/marshaller"
	"github.com/douglas-soares/rpc-quick/src/models"
)

type Invoker interface {
	Register(fnName string, fFunc interface{})
}

type invoker struct {
	funcs map[string]reflect.Value
}

// Register the name of the function and its entries
func (i *invoker) Register(fnName string, fFunc interface{}) {
	if _, ok := i.funcs[fnName]; ok {
		return
	}
	// colocar thread aqui
	i.funcs[fnName] = reflect.ValueOf(fFunc)
}

// fazer funcao de unregister

func (i *invoker) invoke(data []byte) models.Response {
	var req models.Request
	err := marshaller.Unmarshall(data, &req)
	if err != nil {
		return i.returnError(err)
	}

	return i.execute(req)
}

func (i *invoker) execute(req models.Request) models.Response {
	f, ok := i.funcs[req.Function]
	if !ok {
		err := fmt.Errorf("func %s not registered", req.Function)
		return models.Response{Result: nil, Err: err}
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
	resArgs := make([]interface{}, len(out)-1)
	for i := 0; i < len(out)-1; i++ {
		// Interface returns the constant value stored in v as an interface{}.
		resArgs[i] = out[i].Interface()
	}

	// pack error argument
	var err error
	if e, ok := out[len(out)-1].Interface().(error); ok {
		// convert the error into error string value
		err = e
	}
	return models.Response{Result: resArgs, Err: err}
}

func (i *invoker) returnError(err error) models.Response {
	return models.Response{
		Err: err,
	}
}
