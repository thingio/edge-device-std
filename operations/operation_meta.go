package operations

import (
	"encoding/json"
	"github.com/thingio/edge-device-std/msgbus/message"
	"github.com/thingio/edge-device-std/version"
)

type (
	MetaOperationType = OperationType
)

const (
	MetaOperationTypeProductMutation   MetaOperationType = "PRODUCT"
	MetaOperationTypeDeviceMutation    MetaOperationType = "DEVICE"
	MetaOperationTypeDriverInit        MetaOperationType = "INIT"
	MetaOperationTypeDriverHealthCheck MetaOperationType = "STATUS"
)

type MetaOperation struct {
	operation
}

func (o *MetaOperation) Topic() Topic {
	return &commonTopic{
		category: OperationCategoryMeta,
		version:  o.ver,
		tags: map[TopicTagKey]string{
			TopicTagKeyOptMode:    string(o.optMode),
			TopicTagKeyProtocolID: o.protocolID,
			TopicTagKeyOptType:    string(o.optType),
			TopicTagKeyReqID:      o.reqID,
		},
	}
}

func (o *MetaOperation) ToMessage() (*message.Message, error) {
	payload, err := json.Marshal(o.value)
	if err != nil {
		return nil, err
	}
	return &message.Message{
		Topic:   o.Topic().String(),
		Payload: payload,
	}, nil
}

func NewMetaOperation(optMode OperationMode, protocolID string, optType MetaOperationType, reqID string) *MetaOperation {
	o := &MetaOperation{
		operation: operation{
			optCategory: OperationCategoryMeta,
			ver:         version.MetaVersion,
			optMode:     optMode,
			protocolID:  protocolID,
			optType:     optType,
			reqID:       reqID,
		},
	}
	return o
}

func ParseMetaOperation(msg *message.Message) (*MetaOperation, error) {
	topic, err := ParseTopic(msg)
	if err != nil {
		return nil, err
	}
	tags := topic.TagValues()
	o := NewMetaOperation(OperationMode(tags[0]), tags[1], OperationType(tags[2]), tags[3])
	o.payload = msg.Payload
	return o, nil
}
