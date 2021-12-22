package msgbus

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/thingio/edge-device-std/config"
	"github.com/thingio/edge-device-std/logger"
	"github.com/thingio/edge-device-std/msgbus/message"
	"github.com/thingio/edge-device-std/msgbus/mqtt"
)

func NewMessageBus(opts *config.MessageBusOptions, lg *logger.Logger) (MessageBus, error) {
	var mb MessageBus
	switch opts.Type {
	case config.MessageBusTypeMQTT:
		mqttOpts := &opts.MQTT
		if mqttOpts == nil {
			return nil, fmt.Errorf("the configuration for MQTT is required")
		}
		mmb, err := mqtt.NewMQTTMessageBus(mqttOpts, lg)
		if err != nil {
			return nil, err
		}
		mb = mmb
	default:
		return nil, fmt.Errorf("unsupported message bus type: %s", opts.Type)
	}

	if err := mb.Connect(); err != nil {
		return nil, errors.Wrap(err, "fail to connect to the message bus")
	}
	return mb, nil
}

// MessageBus encapsulates all common manipulations based on MQTT.
type MessageBus interface {
	IsConnected() bool

	Connect() error

	Disconnect() error

	Publish(o *message.Message) error

	Subscribe(handler message.Handler, topics ...string) error

	Unsubscribe(topics ...string) error

	Call(request *message.Message, rspTpc, errTpc string) (response *message.Message, err error)
}
