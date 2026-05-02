package errors

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/go-errors/errors"
	"github.com/jgolang/errors/codes"
)

// Error wraps another error with a stack trace, optional friendly message, and optional code.
type Error struct {
	Wrapper *errors.Error
	Message string // A non-technical, user-friendly message describing the error.
	cause   error
	Code    codes.Coder // A custom error code to categorize or identify the error.
}

// StackTraceStr returns the stack trace captured when this Error was created.
func (err *Error) StackTraceStr() string {
	buffer := bytes.NewBufferString("")

	if err == nil || err.Wrapper == nil {
		return "\n\n"
	}

	frames := err.Wrapper.StackFrames()

	if err.Message != "" {
		buffer.WriteString(fmt.Sprintf("\n\n%s\n\n", err.Message))

		buffer.WriteString(fmt.Sprintf("·    Cause: %s\n\n", err.Wrapper.Error()))
	} else {
		buffer.WriteString(fmt.Sprintf("\n\n%s\n\n", err.Wrapper.Error()))
	}

	for _, frame := range frames {
		buffer.WriteString(
			fmt.Sprintf(
				"%s.%s\n",
				frame.Package,
				frame.Name,
			),
		)

		buffer.WriteString(
			fmt.Sprintf(
				"·    %s:%d\n",
				frame.File,
				frame.LineNumber,
			),
		)
	}

	return buffer.String()
}

// StackTrace returns the stack trace as a slog group value.
func (err *Error) StackTrace() slog.Value {
	if err == nil || err.Wrapper == nil {
		return slog.GroupValue()
	}

	frames := err.Wrapper.StackFrames()
	var as []slog.Attr
	for level, frame := range frames {
		fmtFrame := fmt.Sprintf(
			"%s:%d (%s)",
			frame.File,
			frame.LineNumber,
			frame.Name,
		)
		as = append(as, slog.String(fmt.Sprintf("frame_%v", level), fmtFrame))
	}
	return slog.GroupValue(as...)
}

// Unwrap returns the error that caused this Error.
func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}

	return err.cause
}

// Error returns the text of this Error.
func (err *Error) Error() string {
	if err == nil {
		return ""
	}

	var result string

	// Include the custom error code if it exists
	if err.Code != nil {
		result += fmt.Sprintf("[%s]", err.Code.Str())
		if err.Code.Msg() != "" {
			result += fmt.Sprintf("(%s)", err.Code.Msg())
		}
	}

	// Include the custom message if it exists
	if err.Message != "" {
		result += fmt.Sprintf(" %s", err.Message)
	}

	// Always include the original wrapped error message
	if err.Wrapper != nil {
		if result != "" {
			result += ": "
		}
		result += err.Wrapper.Error()
	}

	return result
}
