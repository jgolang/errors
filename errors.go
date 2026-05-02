package errors

import (
	"log/slog"
	"strconv"
	"strings"

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
	if err == nil || err.Wrapper == nil {
		return "\n\n"
	}

	frames := err.Wrapper.StackFrames()
	var buffer strings.Builder

	if err.Message != "" {
		buffer.WriteString("\n\n")
		buffer.WriteString(err.Message)
		buffer.WriteString("\n\n")
		buffer.WriteString("·    Cause: ")
		buffer.WriteString(err.Wrapper.Error())
		buffer.WriteString("\n\n")
	} else {
		buffer.WriteString("\n\n")
		buffer.WriteString(err.Wrapper.Error())
		buffer.WriteString("\n\n")
	}

	for _, frame := range frames {
		buffer.WriteString(frame.Package)
		buffer.WriteString(".")
		buffer.WriteString(frame.Name)
		buffer.WriteString("\n")
		buffer.WriteString("·    ")
		buffer.WriteString(frame.File)
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(frame.LineNumber))
		buffer.WriteString("\n")
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
		fmtFrame := frame.File + ":" + strconv.Itoa(frame.LineNumber) + " (" + frame.Name + ")"
		as = append(as, slog.String("frame_"+strconv.Itoa(level), fmtFrame))
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
		result += "[" + err.Code.Str() + "]"
		if err.Code.Msg() != "" {
			result += "(" + err.Code.Msg() + ")"
		}
	}

	// Include the custom message if it exists
	if err.Message != "" {
		result += " " + err.Message
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
