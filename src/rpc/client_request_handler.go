package rpc

import (
	"context"
	"crypto/tls"
	"fmt"

	quic "github.com/lucas-clemente/quic-go"
)

type requestHandler interface {
	send(location string, msg []byte) (response []byte, err error)
}

type clientRequestHandler struct {
}

// newRequestHandler sada
func newRequestHandler() requestHandler {
	return &clientRequestHandler{}
}

func (h clientRequestHandler) send(location string, msg []byte) ([]byte, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr("localhost:4242", tlsConf, nil)
	if err != nil {
		fmt.Println(1, "client:", err)
	}
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		fmt.Println(2, "client:", err)
	}

	err = send(stream, msg)
	if err != nil {
		fmt.Println(3, "client:", err)
	}

	return read(stream)
}
