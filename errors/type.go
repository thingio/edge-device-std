package errors

import (
	"fmt"
	"net/http"
)

var (
	Types = map[int]ErrType{}

	BadRequest       = NewType(http.StatusBadRequest, "BadRequest")
	NotFound         = NewType(http.StatusNotFound, "BadRequest")
	MethodNotAllowed = NewType(http.StatusMethodNotAllowed, "MethodNotAllowed")
	Internal         = NewType(http.StatusInternalServerError, "Internal")
	Unknown          = NewType(999, "Unknown")
	Configuration    = NewType(100000, "Configuration")
	MessageBus       = NewType(200000, "MessageBus")
	Driver           = NewType(300000, "Driver")
	DeviceTwin       = NewType(400000, "DeviceTwin")
	MetaStore        = NewType(500000, "MetaStore")
	DataStore        = NewType(600000, "DataStore")
)

type ErrType struct {
	Code int    `json:"code"`
	Msg  string `json:"Msg"`
}

func (t *ErrType) Error(format string, args ...interface{}) EdgeError {
	return &CommonEdgeError{
		Msg:           fmt.Sprintf(format, args...),
		wrapped:       nil,
		ErrType:       t,
		stackMessages: getStackMessages(),
	}
}

func (t *ErrType) Cause(err error, format string, args ...interface{}) EdgeError {
	if err == nil {
		return t.Error(format, args...)
	}

	e := Unwrap(err)
	if e.Type().Code == Unknown.Code {
		e.Type().Code = t.Code
	}
	e.stackMessages += getStackMessages()
	return e.addMsg(format, args...)
}

func NewType(code int, msg string) ErrType {
	if _, ok := Types[code]; ok {
		panic(fmt.Errorf("duplicated type code: %d", code))
	}
	tp := ErrType{code, msg}
	Types[code] = tp
	return tp
}
