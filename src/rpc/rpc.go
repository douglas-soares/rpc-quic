package rpc

import (
	"crypto/tls"

	quic "github.com/lucas-clemente/quic-go"
)

type Client interface {
	Call(function string, arg interface{}, result interface{}) error
	Close() error
}

// NewClient creates a new RPC-QUIC client. It enables to connect with the server
func NewClient(serverAddr string, tlsConfig *tls.Config, quicConfig *quic.Config) Client {
	crh := newClientRequestHandler(tlsConfig, quicConfig)
	return newClientProxy(serverAddr, crh)
}

type Server interface {
	Register(function string, fFunc interface{})
	ListenAndServe(addr string, tlsConfig *tls.Config, quicConfig *quic.Config) error
}

type server struct {
	inv *invoker
	srh serverRequestHandler
}

// NewServer creates a new RPC-QUIC server
func NewServer() Server {
	invoker := newInvoker()
	serverRequestHandler := newServerRequestHandler(invoker)

	return &server{
		inv: invoker,
		srh: serverRequestHandler,
	}
}

// Register registers a remote function
func (s *server) Register(function string, fFunc interface{}) {
	s.inv.Register(function, fFunc)
}

// ListenAndServe starts to listen for client connections
func (s *server) ListenAndServe(addr string, tlsConfig *tls.Config, quicConfig *quic.Config) error {
	return s.srh.ListenAndServe(addr, tlsConfig, quicConfig)
}
