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
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{address, timeout, in, out, nil}
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
	_, err := io.Copy(tc.conn, tc.in)
	return err
}

func (tc *TelnetClientImpl) Receive() error {
	_, err := io.Copy(tc.out, tc.conn)
	return err
}
