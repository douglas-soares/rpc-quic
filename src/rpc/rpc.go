package rpc

import (
	"crypto/tls"

	quic "github.com/lucas-clemente/quic-go"
)

// Client con stains
type Client interface {
	Call(result interface{}, function string, args ...interface{}) error
	Close() error
}

// NewClient creates
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

func (s *server) ListenAndServe(addr string, tlsConfig *tls.Config, quicConfig *quic.Config) error {
	return s.srh.ListenAndServe(addr, tlsConfig, quicConfig)
}
