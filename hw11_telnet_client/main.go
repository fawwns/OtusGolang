package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := host + ":" + port

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting: %s\n", err)
		return
	}
	defer client.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := client.Receive(); err != nil {
			fmt.Fprintln(os.Stderr, "...Connection closed by peer")
		}
	}()

	go func() {
		defer wg.Done()
		if err := client.Send(); err != nil {
			fmt.Fprintln(os.Stderr, "...EOF")
		}
	}()

	go func() {
		<-ctx.Done()
		fmt.Fprintln(os.Stderr, "...Connection closed")
		client.Close()
	}()

	wg.Wait()
}
