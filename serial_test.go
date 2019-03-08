package modbus

import (
	"bytes"
	"io"
	"testing"
	"time"
)

type nopCloser struct {
	io.ReadWriter

	closed bool
}

func (n *nopCloser) Close() error {
	n.closed = true
	return nil
}

func TestSerialCloseIdle(t *testing.T) {
	t.Skip()
	port := &nopCloser{
		ReadWriter: &bytes.Buffer{},
	}
	s := serialPort{
		port: port,
	}

	time.Sleep(150 * time.Millisecond)
	if !port.closed || s.port != nil {
		t.Fatalf("serial port is not closed when inactivity: %+v", port)
	}
}
