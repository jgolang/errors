package errors

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/go-errors/errors"
	"github.com/jgolang/errors/codes"
)

// Error is a wrapper of an existing error containing the error stack trace at the moment of creation
type Error struct {
	Wrapper *errors.Error
	Message string // A non-technical, user-friendly message describing the error.
	cause   error
	Code    codes.Coder // A custom error code to categorize or identify the error.
}

// StackTrace Returns an string containing the stack trace computed at the creation moment of this `Error`.
func (err *Error) StackTraceStr() string {
	frames := err.Wrapper.StackFrames()
	buffer := bytes.NewBufferString("")

	if err.Message != "" {
		buffer.WriteString(fmt.Sprintf("\n\n%s\n\n", err.Message))

		buffer.WriteString(fmt.Sprintf("·    Cause: %s\n\n", err.Wrapper.Error()))
	} else if err.Wrapper != nil {
		buffer.WriteString(fmt.Sprintf("\n\n%s\n\n", err.Wrapper.Error()))
	} else {
		buffer.WriteString("\n\n")
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

func (err *Error) StackTrace() slog.Value {
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

// Error returns the text of this `Error`
func (err *Error) Error() string {
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
