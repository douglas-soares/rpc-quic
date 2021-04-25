package rpc

import (
	"crypto/tls"

	quic "github.com/lucas-clemente/quic-go"
)

type clientRequestHandler struct {
	quicConfig *quic.Config
	tlsConfig  *tls.Config
	conn       quic.Stream
	addr       string
}

func newClientRequestHandler(tlsConfig *tls.Config, quicConfig *quic.Config) *clientRequestHandler {
	return &clientRequestHandler{
		quicConfig: quicConfig,
		tlsConfig:  tlsConfig,
	}
}

func (h *clientRequestHandler) dialAndSend(addr string, msg []byte) ([]byte, error) {
	err := h.dial(addr)
	if err != nil {
		return nil, err
	}
	return h.send(msg)
}

func (h *clientRequestHandler) dial(addr string) error {
	if h.conn != nil && h.addr == addr {
		return nil
	}

	h.addr = addr

	session, err := quic.DialAddrEarly(addr, h.tlsConfig, h.quicConfig)
	if err != nil {
		return err
	}

	stream, err := session.OpenStream()
	if err != nil {
		return err
	}

	h.conn = stream

	return nil
}

func (h *clientRequestHandler) send(msg []byte) ([]byte, error) {
	transport := newTransportHelper(h.conn)

	err := transport.send(msg)
	if err != nil {
		if err.Error() == "NO_ERROR: No recent network activity" {
			h.conn = nil
			return h.dialAndSend(h.addr, msg)
		}
		return nil, err
	}

	return transport.read()
}

func (h *clientRequestHandler) close() error {
	err := h.conn.Close()
	h.conn = nil
	return err
}
