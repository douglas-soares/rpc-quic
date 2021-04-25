package rpc

import (
	"bytes"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

func marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func unmarshal(b []byte, v interface{}) error {
	buf := bytes.NewBuffer(b)
	decoder := msgpack.NewDecoder(buf)
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}
