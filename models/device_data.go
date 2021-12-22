package models

import (
	"fmt"
	"reflect"
	"time"
)

type DeviceDataWrapper struct {
	ProductID  string
	DeviceID   string
	FuncID     string
	Properties map[ProductPropertyID]*DeviceData
}

type DeviceData struct {
	Name  string      `json:"name"`  // the name of the data
	Type  string      `json:"type"`  // the type of the raw value
	Value interface{} `json:"value"` // raw value
	Ts    time.Time   `json:"ts"`    // the timestamp of reading the raw value from the real device
}

func NewDeviceData(name string, valueType PropertyValueType, value interface{}) (*DeviceData, error) {
	if err := validate(valueType, value); err != nil {
		return nil, fmt.Errorf("fail to new DeviceData, because %s", err.Error())
	}

	return &DeviceData{
		Name:  name,
		Type:  valueType,
		Value: value,
		Ts:    time.Now(),
	}, nil
}

func (d *DeviceData) String() string {
	return fmt.Sprintf("Data: %s, %s:%s", d.Name, d.Type, d.ValueToString())
}

// ValueToString returns the string format of the Value.
func (d *DeviceData) ValueToString() string {
	return fmt.Sprintf("%v", d.Value)
}

// IntValue returns the Value in string type, and returns errors if the Type is not PropertyValueTypeInt.
func (d *DeviceData) IntValue() (int64, error) {
	var value int64
	if d.Type != PropertyValueTypeInt {
		return value, fmt.Errorf("the expecting type is %s, but the pre-defined type is %s",
			PropertyValueTypeInt, d.Type)
	}

	switch d.Value.(type) {
	case int, int8, int16, int32, int64:
		return d.Value.(int64), nil
	case float32, float64:
		return int64(d.Value.(float64)), nil
	default:
		return value, fmt.Errorf("fail to parse value '%v' using type %s, raw type is %s",
			d.Value, d.Type, reflect.TypeOf(d.Value))
	}
}

// UintValue returns the Value in string type, and returns errors if the Type is not PropertyValueTypeUint.
func (d *DeviceData) UintValue() (uint64, error) {
	var value uint64
	if d.Type != PropertyValueTypeUint {
		return value, fmt.Errorf("the expecting type is %s, but the pre-defined type is %s",
			PropertyValueTypeUint, d.Type)
	}
	switch d.Value.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return d.Value.(uint64), nil
	case float32, float64:
		return uint64(d.Value.(float64)), nil
	default:
		return value, fmt.Errorf("fail to parse value '%v' using type %s, raw type is %s",
			d.Value, d.Type, reflect.TypeOf(d.Value))
	}
}

// FloatValue returns the Value in string type, and returns errors if the Type is not PropertyValueTypeFloat.
func (d *DeviceData) FloatValue() (float64, error) {
	var value float64
	if d.Type != PropertyValueTypeFloat {
		return value, fmt.Errorf("the expecting type is %s, but the pre-defined type is %s",
			PropertyValueTypeFloat, d.Type)
	}

	switch d.Value.(type) {
	case float32, float64:
		return d.Value.(float64), nil
	default:
		return value, fmt.Errorf("fail to parse value '%v' using type %s, raw type is %s",
			d.Value, d.Type, reflect.TypeOf(d.Value))
	}
}

// BoolValue returns the Value in string type, and returns errors if the Type is not PropertyValueTypeBool.
func (d *DeviceData) BoolValue() (bool, error) {
	var value bool
	if d.Type != PropertyValueTypeBool {
		return value, fmt.Errorf("the expecting type is %s, but the pre-defined type is %s",
			PropertyValueTypeBool, d.Type)
	}

	switch d.Value.(type) {
	case bool:
		return d.Value.(bool), nil
	default:
		return value, fmt.Errorf("fail to parse value '%v' using type %s, raw type is %s",
			d.Value, d.Type, reflect.TypeOf(d.Value))
	}
}

// StringValue returns the Value in string type, and returns errors if the Type is not PropertyValueTypeString.
func (d *DeviceData) StringValue() (string, error) {
	var value string
	if d.Type != PropertyValueTypeString {
		return value, fmt.Errorf("the expecting type is %s, but the pre-defined type is %s",
			PropertyValueTypeString, d.Type)
	}

	switch d.Value.(type) {
	case string:
		return d.Value.(string), nil
	default:
		return value, fmt.Errorf("fail to parse value '%v' using type %s, raw type is %s",
			d.Value, d.Type, reflect.TypeOf(d.Value))
	}
}

// validate checks whether value's real type is the given valueType.
func validate(valueType PropertyValueType, value interface{}) error {
	var ok bool
	switch valueType {
	case PropertyValueTypeInt:
		switch value.(type) {
		case int, int8, int16, int32, int64:
			ok = true
		}
	case PropertyValueTypeUint:
		switch value.(type) {
		case uint, uint8, uint16, uint32, uint64:
			ok = true
		}
	case PropertyValueTypeFloat:
		switch value.(type) {
		case float32, float64:
			ok = true
		}
	case PropertyValueTypeBool:
		switch value.(type) {
		case bool:
			ok = true
		}
	case PropertyValueTypeString:
		switch value.(type) {
		case string:
			ok = true
		}
	default:
		return fmt.Errorf("unsupported value's type: %s", valueType)
	}

	if !ok {
		return fmt.Errorf("fail to parse value '%v' using type %s, raw type is %s",
			value, valueType, reflect.TypeOf(value))
	}
	return nil
}
