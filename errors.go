package errors

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/go-errors/errors"
)

// Error is a wrapper of an existing error containing the error stack trace at the moment of creation
type Error struct {
	wrapper *errors.Error
	message string
	cause   error
}

// StackTrace Returns an string containing the stack trace computed at the creation moment of this `Error`.
func (err *Error) StackTraceStr() string {
	frames := err.wrapper.StackFrames()
	buffer := bytes.NewBufferString("")

	if err.message != "" {
		buffer.WriteString(fmt.Sprintf("\n\n%s\n\n", err.message))

		buffer.WriteString(fmt.Sprintf("·    Cause: %s\n\n", err.wrapper.Error()))
	} else if err.wrapper != nil {
		buffer.WriteString(fmt.Sprintf("\n\n%s\n\n", err.wrapper.Error()))
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

func (err *Error) StackTrace() slog.Attr {
	frames := err.wrapper.StackFrames()
	var as []any
	for level, frame := range frames {
		fmtFrame := fmt.Sprintf(
			"%s:%d(%s)",
			frame.File,
			frame.LineNumber,
			frame.Name,
		)
		as = append(as, slog.String(fmt.Sprintf("frame_%v", level), fmtFrame))
	}
	return slog.Group("origin", as...)
}

// Error returns the text of this `Error`
func (err *Error) Error() string {
	if err.message != "" {
		return fmt.Sprintf(err.message)
	}

	if err.wrapper != nil {
		return err.wrapper.Error()
	}

	return ""
}
