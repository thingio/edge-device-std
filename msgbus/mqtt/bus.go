package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thingio/edge-device-std/config"
	"github.com/thingio/edge-device-std/errors"
	"github.com/thingio/edge-device-std/logger"
	"github.com/thingio/edge-device-std/msgbus/message"
	"strconv"
	"time"
)

func NewMQTTMessageBus(opts *config.MQTTMessageBusOptions, lg *logger.Logger) (*MessageBus, errors.EdgeError) {
	mmb := &MessageBus{
		tokenTimeout: time.Millisecond * time.Duration(opts.TokenTimeoutMillisecond),
		callTimeout:  time.Millisecond * time.Duration(opts.MethodCallTimeoutMillisecond),
		qos:          opts.QoS,
		routes:       make(map[string]message.Handler),
		logger:       lg,
	}
	if err := mmb.setClient(opts); err != nil {
		return nil, errors.MessageBus.Cause(err, "fail to initialize the MQTT client")
	}
	return mmb, nil
}

type MessageBus struct {
	client       mqtt.Client
	tokenTimeout time.Duration
	callTimeout  time.Duration
	qos          int

	routes map[string]message.Handler // topic -> handler

	logger *logger.Logger
}

func (mb *MessageBus) IsConnected() bool {
	return mb.client.IsConnected()
}

func (mb *MessageBus) Connect() error {
	if mb.IsConnected() {
		return nil
	}

	token := mb.client.Connect()
	return mb.handleToken(token)
}

func (mb *MessageBus) Disconnect() error {
	if mb.IsConnected() {
		mb.client.Disconnect(2000) // waiting 2s
	}
	return nil
}

func (mb *MessageBus) Publish(msg *message.Message) error {
	mb.logger.Debugf("send message: %s", msg)
	token := mb.client.Publish(msg.Topic, byte(mb.qos), false, msg.Payload)
	return mb.handleToken(token)
}

func (mb *MessageBus) Subscribe(handler message.Handler, topics ...string) error {
	filters := make(map[string]byte)
	for _, topic := range topics {
		mb.routes[topic] = handler
		filters[topic] = byte(mb.qos)
	}
	callback := func(mc mqtt.Client, msg mqtt.Message) {
		go handler(&message.Message{
			Topic:   msg.Topic(),
			Payload: msg.Payload(),
		})
	}

	token := mb.client.SubscribeMultiple(filters, callback)
	return mb.handleToken(token)
}

func (mb *MessageBus) Unsubscribe(topics ...string) error {
	for _, topic := range topics {
		delete(mb.routes, topic)
	}

	token := mb.client.Unsubscribe(topics...)
	return mb.handleToken(token)
}

// Call needs to bind request and response belonging to the same operation,
// otherwise it will cause confusion when multiple operations are executed concurrently.
func (mb *MessageBus) Call(request *message.Message, rspTpc, errTpc string) (response *message.Message, err error) {
	// subscribe response
	ch := make(chan *message.Message)
	if err = mb.Subscribe(func(msg *message.Message) {
		ch <- msg
	}, rspTpc); err != nil {
		return
	}
	errCh := make(chan *message.Message)
	if err = mb.Subscribe(func(msg *message.Message) {
		errCh <- msg
	}, errTpc); err != nil {
		return
	}
	defer func() {
		_ = mb.Unsubscribe(rspTpc)
		close(ch)
		_ = mb.Unsubscribe(errTpc)
		close(errCh)
	}()

	// publish request
	if err = mb.Publish(request); err != nil {
		return
	}
	// waiting for the response
	ticker := time.NewTicker(mb.callTimeout)
	select {
	case msg := <-ch:
		return msg, nil
	case msg := <-errCh:
		return nil, errors.Unmarshal(msg.Payload)
	case <-ticker.C:
		ticker.Stop()
		return nil, errors.MessageBus.Error("call timeout: %dms", mb.tokenTimeout/time.Millisecond)
	}
}

func (mb *MessageBus) handleToken(token mqtt.Token) error {
	if mb.tokenTimeout > 0 {
		token.WaitTimeout(mb.tokenTimeout)
	} else {
		token.Wait()
	}
	if err := token.Error(); err != nil {
		return errors.MessageBus.Cause(token.Error(), "")
	}
	return nil
}

func (mb *MessageBus) setClient(options *config.MQTTMessageBusOptions) error {
	opts := mqtt.NewClientOptions()
	clientID := "edge-device-sub-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	opts.SetClientID(clientID)
	opts.AddBroker(options.GetBroker())
	mb.logger.Infof("the ID of client for the message bus is %s, connecting to %s", clientID, options.GetBroker())
	opts.SetUsername(options.Username)
	opts.SetPassword(options.Password)
	opts.SetConnectTimeout(time.Duration(options.ConnectTimoutMillisecond) * time.Millisecond)
	opts.SetKeepAlive(time.Minute)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(mb.onConnect)
	opts.SetConnectionLostHandler(mb.onConnectLost)
	opts.SetCleanSession(options.CleanSession)

	if options.WithTLS {
		tlsConfig, err := options.NewTLSConfig()
		if err != nil {
			return err
		}
		opts.SetTLSConfig(tlsConfig)
	}

	mb.client = mqtt.NewClient(opts)
	return nil
}

func (mb *MessageBus) onConnect(mc mqtt.Client) {
	reader := mc.OptionsReader()
	mb.logger.Infof("the connection with %s for the message bus has been established.", reader.Servers()[0].String())

	for tpc, hdl := range mb.routes {
		if err := mb.Subscribe(hdl, tpc); err != nil {
			mb.logger.WithError(err).Errorf("fail to resubscribe the topic: %s", tpc)
		}
	}
}

func (mb *MessageBus) onConnectLost(mc mqtt.Client, err error) {
	reader := mc.OptionsReader()
	mb.logger.WithError(err).Errorf("the connection with %s for the message bus has lost, trying to reconnect.",
		reader.Servers()[0].String())
}
