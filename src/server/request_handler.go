package server

type RequestHandler interface {
	ServeAndListen(location string) (err error)
}

type serverRequestHandler struct {
}

// newRequestHandler sada
func newRequestHandler() RequestHandler {
	return &serverRequestHandler{}
}

func (h serverRequestHandler) ServeAndListen(port string) (err error) {
	return nil
}
