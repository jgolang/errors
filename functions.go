package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"

	"github.com/jgolang/errors/codes"
)

// maxStackDepth is the most program counters a stack trace keeps, the same
// limit go-errors applied before v0.8.0.
const maxStackDepth = 50

// callStack holds captured program counters. It is kept behind a pointer, not
// as a slice field, so that Error stays comparable with ==.
type callStack struct {
	pcs   [maxStackDepth]uintptr
	depth int
}

// captureStack returns the program counters of the current goroutine's stack.
// skip is the number of frames to leave out, starting with the function that
// calls captureStack: 1 makes the stack start at the caller of a constructor.
func captureStack(skip int) *callStack {
	stack := &callStack{}
	// 2 leaves out runtime.Callers and captureStack.
	stack.depth = runtime.Callers(2+skip, stack.pcs[:])
	return stack
}

// With wraps err with a formatted friendly message and captures a new stack trace.
// It keeps the code found anywhere in err's chain, or uses codes.Unknown.
func With(err error, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &Error{
		wrapped: err,
		stack:   captureStack(1),
		Message: fmt.Sprintf(message, args...),
		cause:   err,
		Code:    codeOrUnknown(err),
	}
}

// WithC wraps err with a code, formatted friendly message, and new stack trace.
func WithC(err error, code codes.Coder, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &Error{
		wrapped: err,
		stack:   captureStack(1),
		Message: fmt.Sprintf(message, args...),
		cause:   err,
		Code:    code,
	}
}

// Wrap wraps err and captures a stack trace. It keeps the code found anywhere
// in err's chain, or uses codes.Unknown. A sentinel created by Define is
// wrapped with a stack trace that starts at the caller. Any other top-level
// *Error is returned unchanged; one behind another wrapper, such as
// fmt.Errorf, is wrapped.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	if coded, ok := err.(*Error); ok {
		if coded != nil && coded.defined {
			return traceSentinel(coded)
		}
		return coded
	}

	return &Error{
		wrapped: err,
		stack:   captureStack(1),
		cause:   err,
		Code:    codeOrUnknown(err),
	}
}

// Define creates a sentinel error with a code and message but no stack trace,
// for package-level variables compared with errors.Is. Return it through Wrap
// so the stack trace points at the code that failed, not at program start.
func Define(code codes.Coder, message string) error {
	return &Error{Code: code, Message: message, defined: true}
}

// traceSentinel wraps a stackless sentinel with a stack trace that starts at
// the caller of Wrap. The new error reads like one built by NewC and unwraps
// to the sentinel, so errors.Is, errors.As and CodeOf still match it.
func traceSentinel(sentinel *Error) error {
	return &Error{
		wrapped: stderrors.New(sentinel.Message),
		stack:   captureStack(2),
		cause:   sentinel,
		Code:    sentinel.Code,
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
			wrapped: cause,
			stack:   captureStack(1),
			Message: err.Message,
			cause:   err,
			Code:    code,
		}
	}

	return &Error{
		wrapped: err,
		stack:   captureStack(1),
		cause:   err,
		Code:    code,
	}
}

// New creates an error with a formatted message, default code, and stack trace.
func New(format string, args ...interface{}) error {
	return &Error{
		wrapped: fmt.Errorf(format, args...),
		stack:   captureStack(1),
		Code:    codes.Unknown,
	}
}

// NewC creates an error with a code, formatted message, and stack trace.
func NewC(code codes.Coder, format string, args ...interface{}) error {
	return &Error{
		wrapped: fmt.Errorf(format, args...),
		stack:   captureStack(1),
		Code:    code,
	}
}
