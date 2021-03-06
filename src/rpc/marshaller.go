package rpc

import (
	"bytes"
	"encoding/gob"
)

// Marshall be sent over the network.
func marshall(data rpcData) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshall the binary data into the Go struct
func unmarshall(b []byte) (rpcData, error) {
	buf := bytes.NewBuffer(b)
	var data rpcData
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&data); err != nil {
		return rpcData{}, err
	}
	return data, nil
}
