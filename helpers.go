package errors

import (
	stderrors "errors"

	"github.com/jgolang/errors/codes"
)

// CodeOf returns the first jgolang error code found in err's error chain.
func CodeOf(err error) codes.Coder {
	var wrapped *Error
	if stderrors.As(err, &wrapped) {
		return wrapped.Code
	}

	return nil
}

// PublicMessage returns a user-safe message for err without exposing technical details.
func PublicMessage(err error) string {
	code := CodeOf(err)
	if code == nil || code.Msg() == "" {
		return codes.Unknown.Msg()
	}

	return code.Msg()
}
