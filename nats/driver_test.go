package main

import (
	"github.com/nats-io/nats.go"
	"testing"
	"unsafe"
)

func TestRunNATS(t *testing.T) {
	// Mock NATS server
	mockJetStream := &MockJetStream{
		PublishFunc: func(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
			return &nats.PubAck{}, nil
		},
		SubscribeFunc: func(subj string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
			return &nats.Subscription{}, nil
		},
	}

	// Use unsafe.Pointer to cast the mock struct to nats.JetStreamContext
	js := *(*nats.JetStreamContext)(unsafe.Pointer(mockJetStream))

	// Replace the real JetStream context with the mock
	natsConnect := func(url string, opts ...nats.Option) (*nats.Conn, error) {
		return &nats.Conn{}, nil
	}
	*(*func(cfg *nats.StreamConfig) (*nats.StreamInfo, error))(unsafe.Pointer(&js.AddStream)) = func(cfg *nats.StreamConfig) (*nats.StreamInfo, error) { return &nats.StreamInfo{}, nil }
	jetStream := func(nc *nats.Conn, opts ...nats.JSOpt) (nats.JetStreamContext, error) {
		return js, nil
	}

	// Run the function
	if err := runNATS(nil, js); err != nil {
		t.Fatalf("runNATS() failed: %v", err)
	}
}
