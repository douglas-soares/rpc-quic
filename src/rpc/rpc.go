package rpc

import "crypto/tls"

// Client con stains
type Client interface {
	Call(result interface{}, function string, args ...interface{}) error
}

// NewClient creates
func NewClient(serverAddr string, tlsConfig *tls.Config) Client {
	crh := newClientRequestHandler(tlsConfig)
	return newClientProxy(serverAddr, crh)
}

type Server interface {
	Register(function string, fFunc interface{})
	ListenAndServe(addr string, tlsConfig *tls.Config) error
}

type server struct {
	inv *invoker
	srh serverRequestHandler
}

func NewServer() Server {
	invoker := newInvoker()
	serverRequestHandler := newServerRequestHandler(invoker)

	return &server{
		inv: invoker,
		srh: serverRequestHandler,
	}
}

func (s *server) Register(function string, fFunc interface{}) {
	s.inv.Register(function, fFunc)
}

func (s *server) ListenAndServe(addr string, tlsConfig *tls.Config) error {
	return s.srh.ListenAndServe(addr, tlsConfig)
}
