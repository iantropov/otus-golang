package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	closed  bool
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{address, timeout, in, out, nil, false}
}

func (tc *TelnetClientImpl) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	tc.conn = conn
	return err
}

func (tc *TelnetClientImpl) Close() error {
	tc.closed = true
	return tc.conn.Close()
}

func (tc *TelnetClientImpl) Send() error {
	data, err := io.ReadAll(tc.in)
	if err != nil {
		return err
	}
	_, err = tc.conn.Write(data)
	return err
}

func (tc *TelnetClientImpl) Receive() error {
	buf := make([]byte, 512)
	read, err := tc.conn.Read(buf)
	// fmt.Fprintf(os.Stderr, "TELNET - read: %d; len: %d; %v, %s, %v", read, len(buf), buf, buf, err)
	fmt.Fprintf(os.Stderr, "TELNET - read: %d, %v", read, err)
	if read > 0 {
		tc.out.Write(buf[:read])
	}
	return err
}
