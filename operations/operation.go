package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"github.com/thingio/edge-device-std/msgbus/message"
	"github.com/thingio/edge-device-std/version"
)

type (
	OperationCategory string

	OperationType string
	OperationMode string
)

const (
	OperationCategoryMeta OperationCategory = "META"
	OperationCategoryData OperationCategory = "DATA"

	OperationModeUp    OperationMode = "UP"
	OperationModeUpErr OperationMode = "UP-ERR"
	OperationModeDown  OperationMode = "DOWN"
)

type Operation interface {
	Topic() Topic
	ToMessage() (*message.Message, error)

	SetValue(v interface{})
}

type operation struct {
	optCategory OperationCategory
	ver         version.Version
	optMode     OperationMode
	protocolID  string

	optType OperationType
	reqID   string

	value   interface{}
	payload []byte
}

func (o *operation) Topic() Topic {
	return nil
}

func (o *operation) ToMessage() (*message.Message, error) {
	return nil, errors.New("implement me")
}

func (o *operation) SetValue(v interface{}) {
	o.value = v
}

func (o *operation) Unmarshal(v interface{}) error {
	if len(o.payload) == 0 {
		return fmt.Errorf("the payload the operation may not be filled yet")
	}
	return json.Unmarshal(o.payload, v)
}

func NewReqID() string {
	return xid.New().String()
}

func EmptyReqID() string {
	return ""
}
