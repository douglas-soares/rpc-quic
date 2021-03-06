package marshaller

import (
	"bytes"
	"encoding/gob"

	"github.com/douglas-soares/rpc-quick/src/models"
)

// Marshall be sent over the network.
func Marshall(data models.Request) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshall the binary data into the Go struct
func Unmarshall(b []byte) (models.Request, error) {
	buf := bytes.NewBuffer(b)
	var req models.Request
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&req); err != nil {
		return models.Request{}, err
	}
	return req, nil
}
