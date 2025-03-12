package main

import (
	"bufio"
	"net"
	"strings"

	"github.com/fabinzne/go-simple-redis/internal/storage"
)

func startServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	store := storage.NewDataStore()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn, store)
	}
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
  startServer("6379")
}