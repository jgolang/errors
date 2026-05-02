package errors

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jgolang/errors/codes"
)

// With wraps err with a formatted friendly message and captures a new stack trace.
func With(err error, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	var code codes.Coder = codes.Unknown
	if wrapped, ok := err.(*Error); ok {
		if wrapped.Code != nil {
			code = wrapped.Code
		}
	}

	return &Error{
		Wrapper: errors.Wrap(err, 1),
		Message: fmt.Sprintf(message, args...),
		cause:   err,
		Code:    code,
	}
}

// WithC wraps err with a code, formatted friendly message, and new stack trace.
func WithC(err error, code codes.Coder, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &Error{
		Wrapper: errors.Wrap(err, 1),
		Message: fmt.Sprintf(message, args...),
		cause:   err,
		Code:    code,
	}
}

// Wrap wraps err with a default code and captures a stack trace.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*Error); ok {
		return err
	}

	return &Error{
		Wrapper: errors.Wrap(err, 1),
		cause:   err,
		Code:    codes.Unknown,
	}
}

// WrapC wraps err with a code and captures a stack trace.
func WrapC(err error, code codes.Coder) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*Error); ok {
		cause := err.Unwrap()
		if cause == nil {
			cause = err
		}

		return &Error{
			Wrapper: errors.Wrap(cause, 1),
			Message: err.Message,
			cause:   err,
			Code:    code,
		}
	}

	return &Error{
		Wrapper: errors.Wrap(err, 1),
		cause:   err,
		Code:    code,
	}
}

// New creates an error with a formatted message, default code, and stack trace.
func New(format string, args ...interface{}) error {
	return &Error{
		Wrapper: errors.Wrap(fmt.Errorf(format, args...), 1),
		Code:    codes.Unknown,
	}
}

// NewC creates an error with a code, formatted message, and stack trace.
func NewC(code codes.Coder, format string, args ...interface{}) error {
	return &Error{
		Wrapper: errors.Wrap(fmt.Errorf(format, args...), 1),
		Code:    code,
	}
}
