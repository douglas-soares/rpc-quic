package rpc

import (
	"crypto/tls"
	"fmt"

	quic "github.com/lucas-clemente/quic-go"
)

type clientRequestHandler struct {
	quicConfig *quic.Config
	tlsConfig  *tls.Config
	conn       quic.Stream
}

func newClientRequestHandler(tlsConfig *tls.Config, quicConfig *quic.Config) *clientRequestHandler {
	return &clientRequestHandler{
		quicConfig: quicConfig,
		tlsConfig:  tlsConfig,
	}
}

func (h *clientRequestHandler) send(addr string, msg []byte) ([]byte, error) {
	var transport transportHelper
	if h.conn == nil {
		session, err := quic.DialAddrEarly(addr, h.tlsConfig, h.quicConfig)
		if err != nil {
			fmt.Println(1, "client:", err)
		}

		stream, err := session.OpenStream()
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

func (h *clientRequestHandler) close() error {
	err := h.conn.Close()
	h.conn = nil
	return err
}
