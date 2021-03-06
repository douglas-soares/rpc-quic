package rpc

type requestor struct {
	clientRequestHandler
}

// NewRequestor creates a new requestor,
func newRequestor(crh clientRequestHandler) requestor {
	return requestor{
		clientRequestHandler: crh,
	}
}

func (r *requestor) Invoke(location, function string, args []interface{}) interface{} {
	request := rpcData{
		Function: function,
		Args:     args,
	}

	msgMarshalled, err := marshall(request)
	if err != nil {
		return r.returnError(err)
	}

	reqResponse, err := r.send(location, msgMarshalled)
	if err != nil {
		return r.returnError(err)
	}

	response, err := unmarshall(reqResponse)
	if err != nil {
		return r.returnError(err)
	}

	return response
}

func (r *requestor) returnError(err error) rpcData {
	return rpcData{
		Err: err,
	}
}
