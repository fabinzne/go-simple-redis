package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fabinzne/go-simple-redis/internal/storage"
)

var (
	dumpFile = flag.String("d", "./redis.dump", "Filename to dump the data")
	port     = flag.String("p", "6379", "Port to listen on") // Nova flag para porta
	store    *storage.DataStore
)

func StartServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		if isAddrInUse(err) {
			fmt.Printf("Port %s is already in use\n", port)
			os.Exit(1)
		}
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Listening on port %s\n", port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		fmt.Println("\nShutting down...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				fmt.Println("Temporary error when accepting connection - sleeping 1s")
				continue
			}
			break
		}
		go handleConnection(conn, store)
	}
}

func isAddrInUse(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
			return sysErr.Err == syscall.EADDRINUSE
		}
	}
	return false
}

func handleConnection(conn net.Conn, store *storage.DataStore) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)
		if len(parts) < 1 {
			continue
		}

		response := processCommand(parts, store)
		conn.Write([]byte(response + "\n"))
	}
}

func main() {
	flag.Parse()

	store = storage.NewDataStore()

	if err := store.LoadFromFile(*dumpFile); err != nil {
		fmt.Println("Error loading from file:", err)
	}
	
	defer func() {
		fmt.Printf("\nSaving data to %s\n", *dumpFile)
		if err := store.SaveToFile(*dumpFile); err != nil {
			fmt.Println("Error saving to file:", err)
		}
	}()

	StartServer(*port)
}
