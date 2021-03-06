package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/douglas-soares/rpc-quick/src/transport"
	quic "github.com/lucas-clemente/quic-go"
)

type serverRequestHandler struct {
	*invoker
}

// newRequestHandler sada
func NewRequestHandler() serverRequestHandler {
	invoker := newInvoker()
	return serverRequestHandler{invoker: invoker}
}

func (h serverRequestHandler) ServeAndListen(addr string) error {
	listener, err := quic.ListenAddr(addr, GenerateTLSConfig(), nil)
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
			for {
				data, err := transport.Read(stream)
				if err != nil {
					fmt.Println(" error reading mclient msg", err)
					stream.Close() // deverimos fechar?
					return
				}

				response := h.invoker.invoke(data)
				transport.Send(stream, response)
				stream.Close()
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

func GenerateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
