package main

import (
	"bufio"
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
	if tc.scanner.Scan() {
		fmt.Fprintln(os.Stderr, "TELNET:", tc.scanner.Bytes(), tc.scanner.Text(), tc.scanner.Err())
		_, err := tc.out.Write(tc.scanner.Bytes())
		if err != nil {
			return err
		}
		tc.out.Write([]byte{'\n'})
		return nil
	}
	fmt.Fprintln(os.Stderr, "TELNET (AFTER SCAN):", tc.closed, tc.scanner.Err())
	if tc.closed {
		return io.EOF
	}
	err := tc.scanner.Err()
	if err != nil {
		return err
	}
	return io.EOF
}
