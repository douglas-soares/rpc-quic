package rpc

import (
	"fmt"
	"reflect"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

type invoker struct {
	funcs     map[string]reflect.Value
	funcParam map[string]reflect.Type
}

type serverRequest struct {
	Function string              `json:"method"`
	Args     *msgpack.RawMessage `json:"params"`
}

type serverResponse struct {
	Result interface{} `json:"result"`
	Err    string      `json:"error"`
}

func newInvoker() *invoker {
	funcs := make(map[string]reflect.Value)
	funcParam := make(map[string]reflect.Type)
	return &invoker{
		funcs:     funcs,
		funcParam: funcParam,
	}
}

// Register the name of the function and its entries
func (i *invoker) Register(function string, fFunc interface{}) error {
	if fFunc == nil {
		return fmt.Errorf("register error: fFunc is null")
	}

	fFuncType := reflect.TypeOf(fFunc)
	if fFuncType.Kind() != reflect.Func {
		return fmt.Errorf("register error: fFunc must be a function")
	}

	if _, ok := i.funcs[function]; ok {
		return nil
	}

	if fFuncType.NumIn() > 1 {
		return fmt.Errorf("register error: function has more than 1 parameters")
	}

	i.funcs[function] = reflect.ValueOf(fFunc)

	if fFuncType.NumIn() == 1 {
		i.funcParam[function] = fFuncType.In(0)
	}

	fmt.Println("function", function, "registred")

	return nil
}

func (i *invoker) invoke(data []byte) []byte {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in invoke", r)
		}
	}()

	var req serverRequest
	err := unmarshal(data, &req)
	if err != nil {
		fmt.Println("invoker 1", err)
		return i.returnError(err)
	}

	response := i.execute(req)

	marshaledResponse, err := marshal(response)
	if err != nil {
		fmt.Println("invoker 2", err)
		return i.returnError(err)
	}

	return marshaledResponse
}

func (i *invoker) execute(req serverRequest) serverResponse {
	f, ok := i.funcs[req.Function]
	if !ok {
		err := fmt.Sprintf("func %s not registered", req.Function)
		return serverResponse{Result: nil, Err: err}
	}

	paramType, funcHasParam := i.funcParam[req.Function]

	var inArgs []reflect.Value

	if funcHasParam {
		if req.Args == nil {
			err := fmt.Sprintf("func %s requires a parameter", req.Function)
			return serverResponse{Result: nil, Err: err}
		}

		// force pointer
		var argv reflect.Value
		if paramType.Kind() == reflect.Ptr {
			argv = reflect.New(paramType.Elem())
		} else {
			argv = reflect.New(paramType)
		}

		err := unmarshal(*req.Args, argv.Interface())
		if err != nil {
			return serverResponse{Result: nil, Err: err.Error()}
		}

		inArgs = append(inArgs, argv.Elem())
	}

	// invoke requested method
	out := f.Call(inArgs)
	if len(out) == 0 {
		return serverResponse{}
	}

	resArgs := make([]interface{}, len(out))
	for i := 0; i < len(out); i++ {
		resArgs[i] = out[i].Interface()
	}

	return serverResponse{Result: resArgs[0]}
}

func (i *invoker) returnError(err error) []byte {
	resp := serverResponse{
		Err: err.Error(),
	}
	r, _ := marshal(resp)
	return r
}
