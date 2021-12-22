package errors

import (
	"fmt"
	"testing"
)

var (
	RawEmptyError              error = nil
	L1EmptyErrorWithoutMessage       = NewCommonEdgeError(Unknown, "", RawEmptyError)
	L1EmptyErrorWithMessage          = NewCommonEdgeError(Unknown, "with message", RawEmptyError)
	L2EmptyErrorWithoutMessage       = NewCommonEdgeError(Unknown, "", L1EmptyErrorWithoutMessage)
	L2EmptyErrorWithMessage          = NewCommonEdgeError(Unknown, "with message", L1EmptyErrorWithMessage)

	RawError                   = fmt.Errorf("raw errors")
	L1ErrorWithoutMessage      = NewCommonEdgeError(Unknown, "", RawError)
	L1ErrorWithMessage         = NewCommonEdgeError(Unknown, "with message", RawError)
	L2ErrorWithoutMessage      = NewCommonEdgeError(Unknown, "", L1ErrorWithoutMessage)
	L2ErrorWithMessage         = NewCommonEdgeError(Unknown, "with message", L1ErrorWithMessage)
	L2ErrorMixedWithMessage    = NewCommonEdgeError(Unknown, "message", L1ErrorWithoutMessage)
	L3ErrorMixedWithoutMessage = NewCommonEdgeError(Unknown, "", L2ErrorMixedWithMessage)
	L4ErrorMixedWithMessage    = NewCommonEdgeError(Unknown, "message", L3ErrorMixedWithoutMessage)
)

func TestCommonEdgeError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  EdgeError
		want string
	}{
		{"Get all levels of error messages from an empty error without message", L1EmptyErrorWithoutMessage, L1EmptyErrorWithoutMessage.Msg},
		{"Get all levels of error messages from an empty error with message", L1EmptyErrorWithMessage, L1EmptyErrorWithMessage.Msg},
		{"Get all levels of error messages from an EdgeError with 1 empty error wrapped without message", L2EmptyErrorWithoutMessage, L2EmptyErrorWithoutMessage.Msg},
		{"Get all levels of error messages from an EdgeError with 1 empty error wrapped with message", L2EmptyErrorWithMessage,
			fmt.Sprintf("%s\n\t -> %s", L2EmptyErrorWithMessage.Msg, L1EmptyErrorWithMessage.Error())},

		{"Get all levels of error messages from an error without message", L1ErrorWithoutMessage, RawError.Error()},
		{"Get all levels of error messages from an error with message", L1ErrorWithMessage,
			fmt.Sprintf("%s\n\t -> %s", L1ErrorWithMessage.Msg, RawError.Error())},
		{"Get all levels of error message from an EdgeError with 1 error wrapped without message", L2ErrorWithoutMessage, L2ErrorWithoutMessage.Error()},
		{"Get all levels of error message from an EdgeError with 1 error wrapped with message", L2ErrorWithMessage,
			fmt.Sprintf("%s\n\t -> %s", L2ErrorWithMessage.Msg, L1ErrorWithMessage.Error())},
		{"Get all levels of error message from an EdgeError with 1 error wrapped with message fixed", L2ErrorMixedWithMessage,
			fmt.Sprintf("%s\n\t -> %s", L2ErrorMixedWithMessage.Msg, L1ErrorWithoutMessage.Error())},
		{"Get all levels of error message from an EdgeError with 2 error wrapped with message", L3ErrorMixedWithoutMessage, L2ErrorMixedWithMessage.Error()},
		{"Get all levels of error message from an EdgeError with 3 error wrapped with message", L4ErrorMixedWithMessage,
			fmt.Sprintf("%s\n\t -> %s", L4ErrorMixedWithMessage.Msg, L3ErrorMixedWithoutMessage.Error())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommonEdgeError_Message(t *testing.T) {
	tests := []struct {
		name string
		err  EdgeError
		want string
	}{
		{"Get the first level error message from an empty error without message", L1EmptyErrorWithoutMessage, L1EmptyErrorWithoutMessage.Msg},
		{"Get the first level error message from an empty error with message", L1EmptyErrorWithMessage, L1EmptyErrorWithMessage.Msg},
		{"Get the first level error message from an EdgeError with 1 empty error wrapped without message", L2EmptyErrorWithoutMessage, L2EmptyErrorWithoutMessage.Msg},
		{"Get the first level error message from an EdgeError with 1 empty error wrapped with message", L2EmptyErrorWithMessage, L2EmptyErrorWithMessage.Msg},

		{"Get the first level error message from an error without message", L1ErrorWithoutMessage, L1ErrorWithoutMessage.wrapped.Error()},
		{"Get the first level error message from an error with message", L1ErrorWithMessage, L1ErrorWithMessage.Msg},
		{"Get the first level error message from an EdgeError with 1 error wrapped without message", L2ErrorWithoutMessage, L2ErrorWithoutMessage.wrapped.Error()},
		{"Get the first level error message from an EdgeError with 1 error wrapped with message", L2ErrorWithMessage, L2ErrorWithMessage.Msg},
		{"Get the first level error message from an EdgeError with 1 error wrapped with message fixed", L2ErrorMixedWithMessage, L2ErrorMixedWithMessage.Msg},
		{"Get the first level error message from an EdgeError with 2 error wrapped with message", L3ErrorMixedWithoutMessage, L2ErrorMixedWithMessage.Msg},
		{"Get the first level error message from an EdgeError with 3 error wrapped with message", L4ErrorMixedWithMessage, L4ErrorMixedWithMessage.Msg},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Message(); got != tt.want {
				t.Errorf("Msg() = %v, want %v", got, tt.want)
			}
		})
	}
}
