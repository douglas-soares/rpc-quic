package rpc

import (
	"context"
	"crypto/tls"
	"fmt"

	quic "github.com/lucas-clemente/quic-go"
)

type serverRequestHandler struct {
	invoker *invoker
}

// newRequestHandler sada
func newServerRequestHandler(invoker *invoker) serverRequestHandler {
	return serverRequestHandler{invoker: invoker}
}

func (h serverRequestHandler) ListenAndServe(addr string, tlsConfig *tls.Config) error {
	listener, err := quic.ListenAddr(addr, tlsConfig, nil)
	if err != nil {
		return err
	}
	fmt.Println("server started")

	var sess quic.Session
	var stream quic.Stream

	for {
		sess, err = listener.Accept(context.Background())
		if err != nil {
			return err
		}
		go func() {
			stream, err = sess.AcceptStream(context.Background())
			if err != nil {
				fmt.Println(3, "server:", err)
			}
			transport := newTransportHelper(stream)
			for {
				data, err := transport.read()
				if err != nil {
					fmt.Println(" error reading mclient msg", err)
					stream.Close() // deverimos fechar?
					return
				}

				response := h.invoker.invoke(data)
				err = transport.send(response)
				if err != nil {
					fmt.Println(" error sending to client", err)
					stream.Close() // deverimos fechar?
					return
				}
			}
		}()
	}
}

// func (h serverRequestHandler) receiveMessage(conn net.Conn, data interface{}) error {
// 	decoder := gob.NewDecoder(conn) //arrumar isso para n precisar criar
// 	return decoder.Decode(&data)
// }

// func (h serverRequestHandler) sendMessage(conn net.Conn, data interface{}) error {
// 	encoder := gob.NewEncoder(conn) //arrumar isso para n precisar criar
// 	return encoder.Encode(&data)
// }
