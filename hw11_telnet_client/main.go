package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var timeout time.Duration

var ErrInvalidArgs = errors.New("invalid arguments")

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	connectionStr, err := parseConnectionString(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "...Failed to start: %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bufferOut := &bytes.Buffer{}
	bufferIn := &bytes.Buffer{}

	telnetClient := NewTelnetClient(connectionStr, timeout, io.NopCloser(bufferIn), bufferOut)
	err = telnetClient.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "...Failed to connect to: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "...Connected to", connectionStr)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer cancel()
		defer wg.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Fprintln(os.Stderr, "READER:", scanner.Bytes(), scanner.Text(), scanner.Err())
			bufferIn.Reset()
			bufferIn.Write(scanner.Bytes())
			bufferIn.WriteString("\n")
			err := telnetClient.Send()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				fmt.Fprintln(os.Stderr, "...Failed to send data to the telnet client", err)
				return
			}
		}

		if ctx.Err() != nil {
			return
		}

		if scanner.Err() != nil {
			fmt.Fprintln(os.Stderr, "...Received error from scanner", scanner.Err())
			return
		}

		fmt.Fprintln(os.Stderr, "...EOF")
	}()

	go func() {
		defer cancel()
		defer wg.Done()

		for {
			err := telnetClient.Receive()
			fmt.Fprintln(os.Stderr, "WRITER:", bufferOut.Bytes(), err)
			if len(bufferOut.Bytes()) > 0 {
				_, err = os.Stdout.Write(bufferOut.Bytes())
				if err != nil {
					fmt.Fprintln(os.Stderr, "...Failed to send data to STDOUT", err)
					return
				}
				bufferOut.Reset()
			}
			if err != nil {
				if ctx.Err() != nil {
					fmt.Fprintln(os.Stderr, "CONTEXT WAS CLOSED")
					return
				}
				if err == io.EOF {
					fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
					return
				}
				fmt.Fprintln(os.Stderr, "...Failed to receive data from the telnet client", err)
				return
			}
		}
	}()

	<-ctx.Done()
	telnetClient.Close()
	wg.Wait()
}

func parseConnectionString(osArgs []string) (string, error) {
	connectionArgs := make([]string, 0, 2)
	for i := 1; i < len(osArgs); i++ {
		if strings.HasPrefix(osArgs[i], "-") {
			continue
		}
		connectionArgs = append(connectionArgs, osArgs[i])
	}

	if len(connectionArgs) != 2 {
		return "", ErrInvalidArgs
	}

	return fmt.Sprintf("%s:%s", connectionArgs[0], connectionArgs[1]), nil
}
