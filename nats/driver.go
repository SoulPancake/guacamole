package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Call the new function
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	if err := runNATS(nc, js); err != nil {
		log.Fatal(err)
	}
}

func runNATS(_ *nats.Conn, js nats.JetStreamContext) error {

	// Define stream configuration
	streamName := "mystream"
	streamSubjects := []string{"foo"}

	// Add a stream
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: streamSubjects,
	})
	if err != nil {
		return err
	}

	// Publish a message to the stream
	_, err = js.Publish("foo", []byte("Hello, JetStream!"))
	if err != nil {
		return err
	}
	fmt.Println("Message published to stream")

	// Subscribe to the stream
	sub, err := js.SubscribeSync("foo")
	if err != nil {
		return err
	}

	// Fetch a message from the stream
	msg, err := sub.NextMsg(10 * time.Second)
	if err != nil {
		return err
	}
	fmt.Printf("Received message: %s\n", string(msg.Data))

	// Acknowledge the message
	msg.Ack()

	return nil
}
