package rpc

import "crypto/tls"

// ClientRPC con stains
type ClientRPC interface {
	Call(result interface{}, function string, args ...interface{}) error
}

// NewClientRPC creates
func NewClientRPC(serverAddr string, tlsConfig *tls.Config) ClientRPC {
	crh := newClientRequestHandler(tlsConfig)
	return newClientProxy(serverAddr, crh)
}

type ServerRPC interface {
	Register(function string, fFunc interface{})
	ServeAndListen(addr string, tlsConfig *tls.Config) error
}

type serverRPC struct {
	inv *invoker
	srh serverRequestHandler
}

func NewServerRPC() ServerRPC {
	invoker := newInvoker()
	serverRequestHandler := newServerRequestHandler(invoker)

	return &serverRPC{
		inv: invoker,
		srh: serverRequestHandler,
	}
}

func (s *serverRPC) Register(function string, fFunc interface{}) {
	s.inv.Register(function, fFunc)
}

func (s *serverRPC) ServeAndListen(addr string, tlsConfig *tls.Config) error {
	return s.srh.ServeAndListen(addr, tlsConfig)
}
