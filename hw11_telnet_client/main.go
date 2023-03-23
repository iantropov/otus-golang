package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var ErrInvalidArgs = errors.New("invalid arguments")

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	connectionStr, err := parseConnectionString(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "...Failed to start: %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	telnetClient := NewTelnetClient(connectionStr, timeout, os.Stdin, os.Stdout)
	err = telnetClient.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "...Failed to connect to: %v\n", err)
		os.Exit(1)
	}
	defer telnetClient.Close()

	fmt.Fprintln(os.Stderr, "...Connected to", connectionStr)

	go func() {
		defer cancel()
		if err := telnetClient.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "...received error from send: %v\n", err)
		}
	}()

	go func() {
		defer cancel()
		if err := telnetClient.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "...received error from receive: %v\n", err)
		}
	}()

	<-ctx.Done()
}

func parseConnectionString(flagArgs []string) (string, error) {
	if len(flagArgs) != 2 {
		return "", ErrInvalidArgs
	}

	return net.JoinHostPort(flagArgs[0], flagArgs[1]), nil
}
