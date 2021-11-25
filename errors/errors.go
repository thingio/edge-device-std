package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
)

type EdgeError interface {
	// error implements built-in error interface to define a customized error interface.
	error
	Message() string
	StackMessages() string
	Type() *ErrType
}

func NewCommonEdgeError(errType ErrType, msg string, wrapped error) *CommonEdgeError {
	return &CommonEdgeError{
		Msg: msg,
		ErrType: &ErrType{
			Code: errType.Code,
			Msg:  errType.Msg,
		},
		wrapped:       wrapped,
		stackMessages: getStackMessages(),
	}
}

func NewCommonEdgeErrorWrapper(wrapped error) *CommonEdgeError {
	errType := TypeOf(wrapped)
	return &CommonEdgeError{
		Msg:           wrapped.Error(),
		wrapped:       wrapped,
		ErrType:       &errType,
		stackMessages: getStackMessages(),
	}
}

func Unmarshal(data []byte) *CommonEdgeError {
	cee := new(CommonEdgeError)
	if err := json.Unmarshal(data, cee); err != nil {
		return NewCommonEdgeErrorWrapper(err)
	}
	return cee
}

type CommonEdgeError struct {
	// Msg describes the detailed information about the CommonEdgeError.
	Msg string `json:"message"`
	// wrapped is a chain of errors.
	wrapped error
	// ErrType is the type to represent this CommonEdgeError.
	ErrType *ErrType `json:"type"`
	// stackMessages is the information of function call stacks.
	stackMessages string
}

// Error returns all levels of error messages.
func (e *CommonEdgeError) Error() string {
	if e.wrapped == nil {
		return e.Msg
	}

	if e.Msg != "" {
		return e.Msg + "\n\t -> " + e.wrapped.Error()
	} else {
		return e.wrapped.Error()
	}
}

// Message returns the first level error message without further details.
func (e *CommonEdgeError) Message() string {
	if e.Msg == "" && e.wrapped != nil {
		if w, ok := e.wrapped.(*CommonEdgeError); ok {
			return w.Message()
		} else {
			return e.wrapped.Error()
		}
	}

	return e.Msg
}

// StackMessages returns the call stack information of the error function for debugging.
func (e *CommonEdgeError) StackMessages() string {
	if e.wrapped == nil {
		return fmt.Sprintf("%s\n\n%s", e.Msg, e.stackMessages)
	}

	if w, ok := e.wrapped.(*CommonEdgeError); ok {
		return fmt.Sprintf("%s\n\n%s%s", e.Msg, e.stackMessages, w.StackMessages())
	} else {
		return fmt.Sprintf("%s\n\n%s%s", e.Msg, e.stackMessages, e.wrapped.Error())
	}
}

// Type returns the type of this CommonEdgeError.
func (e *CommonEdgeError) Type() *ErrType {
	return e.ErrType
}

// Unwrap returns an errors wrapped in this CommonEdgeError.
func (e *CommonEdgeError) Unwrap() error {
	return e.wrapped
}

func (e *CommonEdgeError) addMsg(format string, args ...interface{}) *CommonEdgeError {
	if format != "" {
		e.Msg = fmt.Sprintf(format, args...) + "\n\t -> " + e.Msg
	}
	return e
}

func Unwrap(err error) *CommonEdgeError {
	if err == nil {
		return nil
	}

	if e, ok := err.(*CommonEdgeError); ok {
		return e
	}
	return Unknown.Error(err.Error()).(*CommonEdgeError)
}

func TypeOf(err error) ErrType {
	var e *CommonEdgeError
	if !errors.As(err, &e) {
		return Unknown
	}
	if e.ErrType.Code != Unknown.Code || e.wrapped == nil {
		return *e.ErrType
	}
	return TypeOf(e.wrapped)
}

func getStackMessages() string {
	pc, filename, line, _ := runtime.Caller(3)
	return fmt.Sprintf("%s\n\t-[%s:%d]\n", runtime.FuncForPC(pc).Name(), filename, line)
}
