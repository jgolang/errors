package errors

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jgolang/errors/codes"
)

// With Creates a new error with a new error message from the `message` supplied as a formatted string plus the `args` parameter. The `err` argument is stored as a reference and a stack trace is computed when this function is called
func With(err error, message string, args ...interface{}) error {
	if wrapped, ok := err.(*Error); ok {
		return &Error{
			Wrapper: wrapped.Wrapper,
			Message: fmt.Sprintf(message, args...),
			cause:   err,
			Code:    codes.Unknown,
		}
	}

	return &Error{
		Wrapper: errors.Wrap(err, 1),
		Message: fmt.Sprintf(message, args...),
		cause:   err,
		Code:    codes.Unknown,
	}
}

// WithC Creates a new error with a new error message from the `message` supplied as a formatted string plus the `args` parameter. The `err` argument is stored as a reference and a stack trace is computed when this function is called
func WithC(err error, code codes.Coder, message string, args ...interface{}) error {
	if wrapped, ok := err.(*Error); ok {
		return &Error{
			Wrapper: wrapped.Wrapper,
			Message: fmt.Sprintf(message, args...),
			cause:   err,
			Code:    code,
		}
	}

	return &Error{
		Wrapper: errors.Wrap(err, 1),
		Message: fmt.Sprintf(message, args...),
		cause:   err,
		Code:    code,
	}
}

// Wrap Creates a new error storing the `err` argument as a reference and a stack trace is computed when this function is called
func Wrap(err error) error {
	if err, ok := err.(*Error); ok {
		return err
	}

	return &Error{
		Wrapper: errors.Wrap(err, 1),
		cause:   err,
		Code:    codes.Unknown,
	}
}

// Wrap Creates a new error storing the `err` argument as a reference and a stack trace is computed when this function is called
func WrapC(err error, code codes.Coder) error {
	if err, ok := err.(*Error); ok {
		return &Error{
			Wrapper: err.Wrapper,
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

// New Creates a new error withouth any reference to another error.  The returned error contains an error message from the `format` supplied as a formatted string plus the `args` parameter.
func New(format string, args ...interface{}) error {
	return &Error{
		Wrapper: errors.Wrap(fmt.Errorf(format, args...), 1),
		Code:    codes.Unknown,
	}
}

// NewError creates a new Error instance with the provided code, message, and wrapper error.
func NewC(code codes.Coder, format string, args ...interface{}) error {
	return &Error{
		Wrapper: errors.Wrap(fmt.Errorf(format, args...), 1),
		Code:    code,
	}
}
