package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

type (
	MessageBusType         = string
	MessageBusProtocolType = string
)

const (
	MessageBusTypeMQTT MessageBusType = "MQTT"

	MessageBusProtocolTypeTCP MessageBusProtocolType = "tcp"
	MessageBusProtocolTypeSSL MessageBusProtocolType = "ssl"
)

type MessageBusOptions struct {
	Type MessageBusType        `json:"type" yaml:"type"`
	MQTT MQTTMessageBusOptions `json:"mqtt" yaml:"mqtt"`
}

type MQTTMessageBusOptions struct {
	// Host is the hostname or IP address of the MQTT broker.
	Host string `json:"host" yaml:"host"`
	// Port is the port of the MQTT broker.
	Port int `json:"port" yaml:"port"`
	// Username is the username of the MQTT broker.
	Username string `json:"username" yaml:"username"`
	// Password is the password of the MQTT broker.
	Password string `json:"password" yaml:"password"`

	// ConnectTimoutMillisecond indicates the timeout of connecting to the MQTT broker.
	ConnectTimoutMillisecond int `json:"connect_timout_millisecond" yaml:"connect_timout_millisecond"`
	// TokenTimeoutMillisecond indicates the timeout of mqtt token.
	TokenTimeoutMillisecond int `json:"token_timeout_millisecond" yaml:"token_timeout_millisecond"`
	// QoS is the abbreviation of MQTT Quality of Service.
	QoS int `json:"qos" yaml:"qos"`
	// CleanSession indicates whether retain messages after reconnecting for QoS1 and QoS2.
	CleanSession bool `json:"clean_session" yaml:"clean_session"`

	// MethodCallTimeoutMillisecond indicates the timeout of method call.
	MethodCallTimeoutMillisecond int `json:"method_call_timeout_millisecond" yaml:"method_call_timeout_millisecond"`

	WithTLS  bool   `json:"with_tls" yaml:"with_tls"`
	CAPath   string `json:"ca_path" yaml:"ca_path"`
	CertPath string `json:"cert_path" yaml:"cert_path"`
	KeyPath  string `json:"key_path" yaml:"key_path"`
}

func (o *MQTTMessageBusOptions) NewTLSConfig() (*tls.Config, error) {
	cert, rootPool, err := o.loadTLSConfig()
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            rootPool,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}, nil
}

func (o *MQTTMessageBusOptions) loadTLSConfig() (tls.Certificate, *x509.CertPool, error) {
	cert, err := tls.LoadX509KeyPair(o.CertPath, o.KeyPath)
	if err != nil {
		return cert, nil, fmt.Errorf("fail to load the certificate files: %s", err.Error())
	}
	rootPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(o.CAPath)
	if err != nil {
		return cert, nil, fmt.Errorf("fail to load the root ca file: %s", err.Error())
	}
	ok := rootPool.AppendCertsFromPEM(caCert)
	if !ok {
		return cert, nil, fmt.Errorf("fail to parse the root ca file")
	}
	return cert, rootPool, nil
}

func (o *MQTTMessageBusOptions) GetBroker() string {
	var protocolType MessageBusProtocolType
	if o.WithTLS {
		protocolType = MessageBusProtocolTypeSSL
	} else {
		protocolType = MessageBusProtocolTypeTCP
	}
	return fmt.Sprintf("%s://%s:%d", protocolType, o.Host, o.Port)
}
