package errors

import (
	"log/slog"
	"runtime"
	"strconv"
	"strings"

	"github.com/jgolang/errors/codes"
)

// Error wraps another error with a stack trace, optional friendly message, and optional code.
// The constructors capture the stack trace of their caller with the runtime
// package; read it with StackTrace or StackTraceStr.
type Error struct {
	Message string // A non-technical, user-friendly message describing the error.
	cause   error
	Code    codes.Coder // A custom error code to categorize or identify the error.

	// wrapped is the error whose text Detail returns. It is set together with
	// stack, so a nil wrapped means the error has no stack trace.
	wrapped error
	// stack holds the program counters captured by the constructor.
	stack *callStack

	// defined marks a sentinel built by Define, which Wrap traces at its call site.
	defined bool
}

// Detail returns the full text of the error this Error wraps, without this
// Error's own code and friendly message: the formatted text for New and NewC,
// the sentinel's message for Wrap of a Define sentinel, and otherwise the text
// of the wrapped error, which includes its own code prefix and message when it
// is an *Error. It returns "" for a nil receiver and for an error without a
// stack trace, such as a Define sentinel.
// Like the stack trace, it can contain implementation details; keep it in
// internal logs and return PublicMessage to users.
func (err *Error) Detail() string {
	if !err.hasStack() {
		return ""
	}

	return err.wrapped.Error()
}

// StackTraceStr returns the stack trace captured when this Error was created.
func (err *Error) StackTraceStr() string {
	if !err.hasStack() {
		return "\n\n"
	}

	frames := err.stackFrames()
	var buffer strings.Builder

	if err.Message != "" {
		buffer.WriteString("\n\n")
		buffer.WriteString(err.Message)
		buffer.WriteString("\n\n")
		buffer.WriteString("·    Cause: ")
		buffer.WriteString(err.Detail())
		buffer.WriteString("\n\n")
	} else {
		buffer.WriteString("\n\n")
		buffer.WriteString(err.Detail())
		buffer.WriteString("\n\n")
	}

	for _, frame := range frames {
		buffer.WriteString(frame.pkg)
		buffer.WriteString(".")
		buffer.WriteString(frame.name)
		buffer.WriteString("\n")
		buffer.WriteString("·    ")
		buffer.WriteString(frame.file)
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(frame.line))
		buffer.WriteString("\n")
	}

	return buffer.String()
}

// StackTrace returns the stack trace as a slog group value.
func (err *Error) StackTrace() slog.Value {
	if !err.hasStack() {
		return slog.GroupValue()
	}

	frames := err.stackFrames()
	var as []slog.Attr
	for level, frame := range frames {
		fmtFrame := frame.file + ":" + strconv.Itoa(frame.line) + " (" + frame.name + ")"
		as = append(as, slog.String("frame_"+strconv.Itoa(level), fmtFrame))
	}
	return slog.GroupValue(as...)
}

// hasStack reports whether err was built by a constructor that captured a
// stack trace; Define sentinels and hand-built values have none.
func (err *Error) hasStack() bool {
	return err != nil && err.wrapped != nil
}

// stackFrame is one resolved frame of a captured stack trace.
type stackFrame struct {
	pkg  string
	name string
	file string
	line int
}

// stackFrames resolves the captured program counters into frames on every
// call. Nothing is cached on the shared *Error, so concurrent reads are safe.
func (err *Error) stackFrames() []stackFrame {
	pcs := err.stack.pcs[:err.stack.depth]
	frames := make([]stackFrame, 0, len(pcs))
	callersFrames := runtime.CallersFrames(pcs)
	for more := len(pcs) > 0; more; {
		var frame runtime.Frame
		frame, more = callersFrames.Next()
		pkg, name := splitFunctionName(frame.Function)
		frames = append(frames, stackFrame{pkg: pkg, name: name, file: frame.File, line: frame.Line})
	}
	return frames
}

// splitFunctionName splits a fully qualified function name into its package
// path and its name, as go-errors did before v0.8.0 so the printed frames keep
// their format: the package ends at the first "." after the last "/".
func splitFunctionName(function string) (pkg, name string) {
	name = function
	if lastSlash := strings.LastIndex(name, "/"); lastSlash >= 0 {
		pkg = name[:lastSlash+1]
		name = name[lastSlash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}
	return pkg, strings.ReplaceAll(name, "·", ".")
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

	if err.hasStack() {
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
	full := err.Detail()
	inner, ok := err.wrapped.(*Error)
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
