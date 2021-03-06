package rpc

import (
	"encoding/binary"
	"io"

	quic "github.com/lucas-clemente/quic-go"
)

type transportHelper struct {
	conn quic.Stream
}

func newTransportHelper(conn quic.Stream) transportHelper {
	return transportHelper{
		conn: conn,
	}
}

func (t transportHelper) send(data []byte) error {
	// we will need 4 more byte then the len of data
	// as TLV header is 4bytes and in this header
	// we will encode how much byte of data
	// we are sending for this request.
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	_, err := t.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// Read TLV sent over the wire
func (t transportHelper) read() ([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(header)

	data := make([]byte, dataLen)
	_, err = io.ReadFull(t.conn, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
