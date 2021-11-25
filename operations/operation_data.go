package operations

import (
	"encoding/json"
	"github.com/thingio/edge-device-std/models"
	"github.com/thingio/edge-device-std/msgbus/message"
	"github.com/thingio/edge-device-std/version"
)

type (
	DataOperationType        = OperationType // the type of device data's operation
	DevicePropertyReportMode = string        // the mode of device property's reporting
)

const (
	DataOperationTypeHealthCheck DataOperationType = "STATUS"    // Device Health Check
	DataOperationTypeRead        DataOperationType = "READ"      // Device Property Soft Read
	DataOperationTypeHardRead    DataOperationType = "HARD-READ" // Device Property Hard Read
	DataOperationTypeWrite       DataOperationType = "WRITE"     // Device Property Write
	DataOperationTypeWatch       DataOperationType = "PROPS"     // Device Property Watch
	DataOperationTypeEvent       DataOperationType = "EVENT"     // Device Event
	DataOperationTypeCall        DataOperationType = "CALL"      // Device Method

	DeviceDataReportModePeriodical DevicePropertyReportMode = "periodical" // report device data at intervals, e.g. 5s, 1m, 0.5h
	DeviceDataReportModeOnChange   DevicePropertyReportMode = "onchange"   // report device data on change
)

type DataOperation struct {
	operation

	productID string
	deviceID  string
	funcID    models.ProductFuncID
}

func (o *DataOperation) Topic() Topic {
	return &commonTopic{
		category: OperationCategoryData,
		version:  o.ver,
		tags: map[TopicTagKey]string{
			TopicTagKeyOptMode:    string(o.optMode),
			TopicTagKeyProtocolID: o.protocolID,
			TopicTagKeyProductID:  o.productID,
			TopicTagKeyDeviceID:   o.deviceID,
			TopicTagKeyFuncID:     o.funcID,
			TopicTagKeyOptType:    string(o.optType),
			TopicTagKeyReqID:      o.reqID,
		},
	}
}

func (o *DataOperation) ToMessage() (*message.Message, error) {
	payload, err := json.Marshal(o.value)
	if err != nil {
		return nil, err
	}

	return &message.Message{
		Topic:   o.Topic().String(),
		Payload: payload,
	}, nil
}

func NewDataOperation(optMode OperationMode, protocolID, productID, deviceID string, funcID models.ProductFuncID,
	optType DataOperationType, reqID string) *DataOperation {
	return &DataOperation{
		operation: operation{
			optCategory: OperationCategoryData,
			ver:         version.DataVersion,
			optMode:     optMode,
			protocolID:  protocolID,
			optType:     optType,
			reqID:       reqID,
		},
		productID: productID,
		deviceID:  deviceID,
		funcID:    funcID,
	}
}

func ParseDataOperation(msg *message.Message) (*DataOperation, error) {
	topic, err := ParseTopic(msg)
	if err != nil {
		return nil, err
	}
	tags := topic.TagValues()
	o := NewDataOperation(OperationMode(tags[0]), tags[1], tags[2], tags[3], tags[4], OperationType(tags[5]), tags[6])
	o.payload = msg.Payload
	return o, nil
}
