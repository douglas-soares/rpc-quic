package marshaller

import (
	"bytes"
	"encoding/gob"
)

// Marshall be sent over the network.
func Marshall(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshall the binary data into the Go struct
func Unmarshall(b []byte, data interface{}) error {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	return nil
}
