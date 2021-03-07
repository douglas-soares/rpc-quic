package rpc

import (
	"context"
	"crypto/tls"
	"fmt"

	quic "github.com/lucas-clemente/quic-go"
)

type clientRequestHandler struct {
	tlsConfig *tls.Config
	conn      quic.Stream
}

func newClientRequestHandler(tlsConfig *tls.Config) *clientRequestHandler {
	return &clientRequestHandler{
		tlsConfig: tlsConfig,
	}
}

func (h *clientRequestHandler) send(addr string, msg []byte) ([]byte, error) {
	var transport transportHelper
	if h.conn == nil {

		session, err := quic.DialAddr(addr, h.tlsConfig, nil)
		if err != nil {
			fmt.Println(1, "client:", err)
		}

		stream, err := session.OpenStreamSync(context.Background())
		if err != nil {
			fmt.Println(2, "client:", err)
		}

		h.conn = stream
	}

	transport = newTransportHelper(h.conn)

	err := transport.send(msg)
	if err != nil {
		if err.Error() == "NO_ERROR: No recent network activity" {
			h.conn = nil
			return h.send(addr, msg)
		}
	}

	return transport.read()
}
