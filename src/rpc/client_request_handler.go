package rpc

import (
	"context"
	"crypto/tls"
	"fmt"

	quic "github.com/lucas-clemente/quic-go"
)

type clientRequestHandler struct {
	tlsConfig *tls.Config
}

func newClientRequestHandler(tlsConfig *tls.Config) clientRequestHandler {
	return clientRequestHandler{
		tlsConfig: tlsConfig,
	}
}

func (h *clientRequestHandler) send(addr string, msg []byte) ([]byte, error) {
	fmt.Println(addr)
	session, err := quic.DialAddr(addr, h.tlsConfig, nil)
	if err != nil {
		fmt.Println(1, "client:", err)
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		fmt.Println(2, "client:", err)
	}
	transport := newTransportHelper(stream)

	err = transport.send(msg)
	if err != nil {
		fmt.Println(3, "client:", err)
	}

	return transport.read()
}
