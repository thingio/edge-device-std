package message

import (
	"encoding/json"
	"fmt"
)

type Handler func(msg *Message)

// Message is an intermediate data format between MQ and MessageBus.
type Message struct {
	Topic   string
	Payload []byte
}

func (m *Message) String() string {
	return fmt.Sprintf("%s: %dbytes", m.Topic, len(m.Payload))
}

func (m *Message) Unmarshal(v interface{}) error {
	return json.Unmarshal(m.Payload, v)
}
