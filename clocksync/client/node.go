package main

import (
	"fmt"
	"net"
	"time"
)

const (
	serverAddress = "localhost:9000"
)

func synchronizeClock() {
	// Connect to the server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Get the current time before sending the request (for round-trip calculation)
	clientTimeBeforeRequest := time.Now()

	// Send a request to the server (client doesn't send anything, just a connection)
	_, err = conn.Write([]byte("sync"))
	if err != nil {
		fmt.Println("Error sending data to server:", err)
		return
	}

	// Receive the server's time
	var serverTimeStr string
	_, err = fmt.Fscan(conn, &serverTimeStr)
	if err != nil {
		fmt.Println("Error receiving data from server:", err)
		return
	}

	// Parse the server's time in RFC3339 format
	serverTime, err := time.Parse(time.RFC3339, serverTimeStr)
	if err != nil {
		fmt.Println("Error parsing server time:", err)
		return
	}

	// Get the current time after receiving the server time (for round-trip calculation)
	clientTimeAfterRequest := time.Now()

	// Calculate the round-trip delay (half it for one-way delay)
	roundTripDelay := clientTimeAfterRequest.Sub(clientTimeBeforeRequest) / 2

	// Adjust the client's clock by factoring in the round-trip delay
	adjustedClientTime := serverTime.Add(roundTripDelay)

	// Calculate the drift (the difference between adjusted client time and server time)
	drift := adjustedClientTime.Sub(serverTime)

	// Print the results
	fmt.Println("Server Time:", serverTime)
	fmt.Println("Adjusted Client Time:", adjustedClientTime)
	fmt.Println("Drift (Adjusted Client Time - Server Time):", drift)
}

func main() {
	// Synchronize the clock with the server
	synchronizeClock()
}
