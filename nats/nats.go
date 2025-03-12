package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type MockJetStream struct {
	PublishFunc   func(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error)
	SubscribeFunc func(subj string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error)
}

func (m *MockJetStream) Publish(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	if m.PublishFunc != nil {
		return m.PublishFunc(subj, data, opts...)
	}
	return nil, fmt.Errorf("Publish method not implemented")
}

func (m *MockJetStream) Subscribe(subj string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	if m.SubscribeFunc != nil {
		return m.SubscribeFunc(subj, cb, opts...)
	}
	return nil, fmt.Errorf("Subscribe method not implemented")
}

//func main() {
//	mockJetStream := &MockJetStream{
//		PublishFunc: func(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
//			fmt.Println("Mock Publish called")
//			return &nats.PubAck{}, nil
//		},
//		SubscribeFunc: func(subj string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
//			fmt.Println("Mock Subscribe called")
//			return &nats.Subscription{}, nil
//		},
//	}
//
//	// Use unsafe.Pointer to cast the mock struct to nats.JetStreamContext
//	js := *(*nats.JetStreamContext)(unsafe.Pointer(mockJetStream))
//
//	// Now you can use js as if it were a nats.JetStreamContext
//	js.Publish("test.subject", []byte("test message"))
//	js.Subscribe("test.subject", func(msg *nats.Msg) {})
//}
