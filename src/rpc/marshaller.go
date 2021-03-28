package rpc

import (
	"bytes"
	"encoding/json"
)

// marshal be sent over the network.
func marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal the binary data into the Go struct
func unmarshal(b []byte, v interface{}) error {
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}
