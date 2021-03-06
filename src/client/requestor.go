package client

import (
	"github.com/douglas-soares/rpc-quick/src/marshaller"
	"github.com/douglas-soares/rpc-quick/src/models"
)

// Requestor contains
type Requestor interface {
	Invoke(location, function string, args []interface{}) interface{}
}

type requestor struct {
	requestHandler
}

// NewRequestor creates a new requestor,
func NewRequestor(h requestHandler) Requestor {
	return &requestor{
		requestHandler: h,
	}
}

func (r *requestor) Invoke(location, function string, args []interface{}) interface{} {
	request := models.Request{
		Function: function,
		Args:     args,
	}

	msgMarshalled, err := marshaller.Marshall(request)
	if err != nil {
		return r.returnError(err)
	}

	reqResponse, err := r.send(location, msgMarshalled)
	if err != nil {
		return r.returnError(err)
	}

	response, err := marshaller.Unmarshall(reqResponse)
	if err != nil {
		return r.returnError(err)
	}

	return response
}

func (r *requestor) returnError(err error) models.Request {
	return models.Request{
		Err: err,
	}
}
