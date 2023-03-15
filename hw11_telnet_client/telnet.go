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
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return &TelnetClientImpl{address, timeout, in, out, nil}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.

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
	scanner := bufio.NewScanner(tc.conn)
	ok := scanner.Scan()
	if ok {
		_, err := tc.out.Write(scanner.Bytes())
		return err
	}
	return scanner.Err()
}
