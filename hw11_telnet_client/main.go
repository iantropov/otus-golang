package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

type MyBuffer struct {
	bytes.Buffer
	onClose func()
}

func NewMyBuffer(onClose func()) *MyBuffer {
	return &MyBuffer{
		onClose: onClose,
	}
}

func (mb *MyBuffer) Close() error {
	mb.onClose()
	return nil
}

func main() {
	// if len(os.Args) != 3 {
	// 	fmt.Println("Please, provide host and port arguments")
	// 	return
	// }

	// host := os.Args[1]
	// port := os.Args[2]

	host := "localhost"
	port := "9000"

	connectionStr := fmt.Sprintf("%s:%s", host, port)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bufferIn := NewMyBuffer(func() {
		fmt.Println("Received close from the server!")
		cancel()
	})
	bufferOut := &bytes.Buffer{}

	telnetClient := NewTelnetClient(connectionStr, timeout, bufferIn, bufferOut)
	err := telnetClient.Connect()
	if err != nil {
		log.Fatalf("failed to connect to: %v\n", err)
	}
	defer telnetClient.Close()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Ready to Scan!")
		for scanner.Scan() {
			if ctx.Err() != nil {
				fmt.Println("Received ctx.Done() in sender")
				return
			}
			fmt.Println("Scanned!")

			bufferIn.Reset()
			bufferIn.Write(scanner.Bytes())
			bufferIn.WriteString("\n")
			err := telnetClient.Send()
			if err != nil {
				fmt.Println("Failed to send data to the telnet client", err)
				return
			}
		}
		fmt.Println("Finished scanning!")

		if scanner.Err() != io.EOF {
			fmt.Println("Received error from scanner", scanner.Err())
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()

		fmt.Println("Ready to receive!")
		for {
			if ctx.Err() != nil {
				fmt.Println("Received ctx.Done() in receiver")
				return
			}
			fmt.Println("Will receive!")

			err := telnetClient.Receive()
			fmt.Print("Received: ")
			if err != nil {
				fmt.Println("Failed to receive data from the telnet client", err)
				return
			}
			_, err = os.Stdout.Write(bufferOut.Bytes())
			if err != nil {
				fmt.Println("Failed to send data to STDOUT", err)
			}
			os.Stdout.WriteString("\n")
			bufferOut.Reset()
		}
	}()

	wg.Wait()

	fmt.Println("FINISHED")
}
