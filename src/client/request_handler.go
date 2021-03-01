package client

type requestHandler interface {
	send(location string, msg []byte) (response []byte, err error)
}

type clientRequestHandler struct {
}

// newRequestHandler sada
func newRequestHandler() requestHandler {
	return &clientRequestHandler{}
}

func (h clientRequestHandler) send(location string, msg []byte) (response []byte, err error) {
	return nil, nil
}
