package errors

import (
	stderrors "errors"

	"github.com/jgolang/errors/codes"
)

// CodeOf returns the first non-nil jgolang error code found in err's error
// chain, including codes behind other wrappers such as fmt.Errorf.
func CodeOf(err error) codes.Coder {
	for err != nil {
		var wrapped *Error
		// errors.As also matches a typed-nil *Error, which has no code or cause.
		if !stderrors.As(err, &wrapped) || wrapped == nil {
			return nil
		}
		if wrapped.Code != nil {
			return wrapped.Code
		}
		// A nil code must not hide a code deeper in the chain.
		err = wrapped.cause
	}

	return nil
}

// KindOf returns the kind declared by CodeOf(err), or an empty Kind when that
// code declares none or err has no code. Reading the kind from the same code
// as CodeOf keeps a status derived from it consistent with the reported code.
func KindOf(err error) codes.Kind {
	if kinded, ok := CodeOf(err).(codes.Kinded); ok {
		return kinded.Kind()
	}

	return ""
}

// codeOrUnknown returns the code found in err's chain, or codes.Unknown.
func codeOrUnknown(err error) codes.Coder {
	if code := CodeOf(err); code != nil {
		return code
	}

	return codes.Unknown
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// PublicMessage returns a user-safe message for err without exposing technical details.
func PublicMessage(err error) string {
	code := CodeOf(err)
	if code == nil || code.Msg() == "" {
		return codes.Unknown.Msg()
	}

	return code.Msg()
}
