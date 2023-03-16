package main

import (
	"bufio"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
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
	scanner *bufio.Scanner
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{address, timeout, in, out, nil, false, nil}
}

func (tc *TelnetClientImpl) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	tc.conn = conn
	if err == nil {
		tc.scanner = bufio.NewScanner(tc.conn)
		return nil
	}
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
	var err error
	var wrote int
	if tc.scanner.Scan() {
		wrote, err = tc.out.Write(tc.scanner.Bytes())
		tc.out.Write([]byte{'\n'})
	}
	if err == nil {
		if wrote == 0 {
			return io.EOF
		}
		err = tc.scanner.Err()
	}
	if err != nil {
		if tc.closed {
			return io.EOF
		}
		return err
	}
	return nil
}
