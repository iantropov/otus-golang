package main

import (
	"io"
	"net"
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
	buffer  []byte
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	buf := make([]byte, 512)
	return &TelnetClientImpl{address, timeout, in, out, nil, buf}
}

func (tc *TelnetClientImpl) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	tc.conn = conn
	return err
}

func (tc *TelnetClientImpl) Close() error {
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
	read, err := tc.conn.Read(tc.buffer)
	if read > 0 {
		tc.out.Write(tc.buffer[:read])
	}
	return err
}
