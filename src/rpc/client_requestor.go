package rpc

import msgpack "github.com/vmihailenco/msgpack/v5"

type requestor struct {
	requestHandler *clientRequestHandler
}

type clientRequest struct {
	Function string      `json:"method"`
	Args     interface{} `json:"params"`
}

type clientResponse struct {
	Result *msgpack.RawMessage `json:"result"`
	Err    string              `json:"error"`
}

// NewRequestor creates a new requestor,
func newRequestor(crh *clientRequestHandler) requestor {
	return requestor{
		requestHandler: crh,
	}
}

func (r *requestor) invoke(location, function string, args interface{}) clientResponse {
	request := clientRequest{
		Function: function,
		Args:     args,
	}

	message, err := marshal(request)
	if err != nil {
		return r.returnError(err)
	}

	reqResponse, err := r.requestHandler.dialAndSend(location, message)
	if err != nil {
		return r.returnError(err)
	}

	var response clientResponse
	err = unmarshal(reqResponse, &response)
	if err != nil {
		return r.returnError(err)
	}

	return response
}

func (r *requestor) Close() error {
	return r.requestHandler.close()
}

func (r *requestor) returnError(err error) clientResponse {
	return clientResponse{
		Err: err.Error(),
	}
}
