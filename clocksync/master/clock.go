package main

import (
	"fmt"
	"net"
	"time"
)

const (
	port = ":9000"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Get the current server time (master clock)
	currentTime := time.Now().Format(time.RFC3339) // Use RFC3339 format

	// Send the current time to the client in RFC3339 format
	_, err := conn.Write([]byte(currentTime))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
	}
}

func startServer() {
	// Listen for incoming connections
	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listen.Close()

	fmt.Println("Server started on", port)

	// Accept connections from clients
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the client in a separate goroutine
		go handleClient(conn)
	}
}

func main() {
	// Start the master clock server
	startServer()
}
