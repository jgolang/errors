package errors

import (
	"fmt"

	"github.com/go-errors/errors"
)

// With Creates a new error with a new error message from the `message` supplied as a formatted string plus the `args` parameter. The `err` argument is stored as a reference and a stack trace is computed when this function is called
func With(err error, message string, args ...interface{}) error {
	if wrapped, ok := err.(*Error); ok {
		return &Error{
			wrapper: wrapped.wrapper,
			message: fmt.Sprintf(message, args...),
			cause:   err,
		}
	}

	return &Error{
		wrapper: errors.Wrap(err, 1),
		message: fmt.Sprintf(message, args...),
		cause:   err,
	}
}

// Wrap Creates a new error storing the `err` argument as a reference and a stack trace is computed when this function is called
func Wrap(err error) error {
	if err, ok := err.(*Error); ok {
		return err
	}

	return &Error{
		wrapper: errors.Wrap(err, 1),
		cause:   err,
	}
}

// New Creates a new error withouth any reference to another error.  The returned error contains an error message from the `format` supplied as a formatted string plus the `args` parameter.
func New(format string, args ...interface{}) error {
	return &Error{
		wrapper: errors.Wrap(fmt.Errorf(format, args...), 1),
	}
}
