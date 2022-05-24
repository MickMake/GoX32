package gosc

import (
	"bufio"
	"bytes"
	"net"
)

// Transport interface describes the transportation used by the Client.
type Transport interface {
	Send(pack Package) error
	Receive() (pack Package, err error)
}

type transport struct {
	conn       net.Conn
	bufferSize int
}

// NewUDPTransport returns the default UDP Transport.
func NewUDPTransport(address string, bufferSize int) (Transport, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}
	return &transport{
		conn:       conn,
		bufferSize: bufferSize,
	}, nil
}

// Send uses buffering to send a complete Package on the UDP socket.
func (t *transport) Send(pack Package) error {
	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)

	err := writePackage(pack, w)
	if err != nil {
		return err
	}
	_ = w.Flush()
	// fmt.Printf("X32 OSC: %v\n", buf.String())	// @TODO - DEBUG
	// fmt.Printf("X32 OSC: %v\n", buf.Bytes())		// @TODO - DEBUG
	_, err = t.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Receive reads the UDP socket buffered and returns a Package when found.
func (t *transport) Receive() (pack Package, err error) {
	buf := make([]byte, t.bufferSize)
	_, err = t.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReaderSize(bytes.NewReader(buf), len(buf))
	return readPackage(r)
}
