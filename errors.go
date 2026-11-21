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

	// defined marks a sentinel built by Define, which Wrap traces at its call site.
	defined bool
}

// StackTraceStr returns the stack trace captured when this Error was created.
func (err *Error) StackTraceStr() string {
	if err == nil || err.Wrapper == nil {
		return "\n\n"
	}

	frames := err.stackFrames()
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

	frames := err.stackFrames()
	var as []slog.Attr
	for level, frame := range frames {
		fmtFrame := frame.File + ":" + strconv.Itoa(frame.LineNumber) + " (" + frame.Name + ")"
		as = append(as, slog.String("frame_"+strconv.Itoa(level), fmtFrame))
	}
	return slog.GroupValue(as...)
}

// stackFrames resolves the captured program counters into frames on every call.
// go-errors' StackFrames caches them on the shared Wrapper without
// synchronization, so reading it from several goroutines is a data race.
// NewStackFrame keeps the frame text identical to StackFrames.
func (err *Error) stackFrames() []errors.StackFrame {
	callers := err.Wrapper.Callers()
	frames := make([]errors.StackFrame, len(callers))
	for i, pc := range callers {
		frames[i] = errors.NewStackFrame(pc)
	}
	return frames
}

// Unwrap returns the error that caused this Error.
func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}

	return err.cause
}

// Error returns the text of this Error: the code and its message in brackets,
// then the friendly message and the wrapped error's text. A friendly message
// or wrapped text equal to the code's message is printed only once. When the
// error directly wraps an *Error with the same code (With, WithC and WrapC over
// a coded error), the wrapped text leaves out its repeated [code](message)
// prefix; behind another wrapper such as fmt.Errorf it is printed as is.
func (err *Error) Error() string {
	if err == nil {
		return ""
	}

	result := err.prefix()
	codeText := ""
	if err.Code != nil {
		codeText = err.Code.Msg()
	}

	// Domain errors often reuse the code's text as their message (ERR-05);
	// printing it again adds noise to every log line.
	if err.Message != "" && err.Message != codeText {
		result += " " + err.Message
	}

	if err.Wrapper != nil {
		wrapped, onlyPrefix := err.wrappedText()
		if !onlyPrefix && (codeText == "" || wrapped != codeText) {
			if result != "" {
				result += ": "
			}
			result += wrapped
		}
	}

	return result
}

// prefix returns "[code](message)", "[code]" when the code has no message,
// or "" for a nil code.
func (err *Error) prefix() string {
	if err.Code == nil {
		return ""
	}
	if err.Code.Msg() == "" {
		return "[" + err.Code.Str() + "]"
	}
	return "[" + err.Code.Str() + "](" + err.Code.Msg() + ")"
}

// wrappedText returns the text of the wrapped error. When it is an *Error with
// the same code, its own prefix and the separator after it are cut from its
// text, and onlyPrefix reports that nothing else was left to print.
func (err *Error) wrappedText() (text string, onlyPrefix bool) {
	full := err.Wrapper.Error()
	inner, ok := err.Wrapper.Err.(*Error)
	if !ok || inner == nil || !sameCode(inner.Code, err.Code) {
		return full, false
	}

	rest := strings.TrimPrefix(full, inner.prefix())
	if rest == "" {
		return "", true
	}
	switch {
	case strings.HasPrefix(rest, ": "):
		rest = rest[len(": "):]
	case strings.HasPrefix(rest, " "):
		rest = rest[len(" "):]
	}
	if rest == "" {
		// The inner error ended with an empty wrapped text; keep its marker.
		return full, false
	}
	return rest, false
}

// sameCode reports whether a and b have the same code string and message.
func sameCode(a, b codes.Coder) bool {
	return a != nil && b != nil && a.Str() == b.Str() && a.Msg() == b.Msg()
}
