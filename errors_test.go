package errors

import (
	stderrors "errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/jgolang/errors/codes"
)

func TestNew(t *testing.T) {
	// Test case parameters
	format := "This is a test message: %s"
	arg := "testArg"

	// Create a new error using the New function
	err := New(format, arg)

	// Assertions
	if err == nil {
		t.Errorf("Expected error to be non-nil")
		return
	}

	// Type assertion
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("Expected error to be of type *Error, got %T", err)
		return
	}

	// Check if the code is set to ErrGenUnknown
	if e.Code != codes.Unknown {
		t.Errorf("Expected error code to be ErrGenUnknown, got %v", e.Code)
	}

	// Check if the error message is formatted correctly
	expectedMsg := fmt.Sprintf("[%s](%s): %s", codes.Unknown.Str(), codes.Unknown.Msg(), fmt.Sprintf(format, arg))
	if e.Error() != expectedMsg {
		t.Errorf("Expected error message to be '%s', got '%s'", expectedMsg, e.Error())
	}
}

func TestNewError(t *testing.T) {
	// Test case parameters
	code := codes.AppValidation // Example code, can be replaced with any code
	format := "Validation failed for field: %s"
	arg := "username"

	Wrap(nil)

	// Create a new error using the NewError function
	err := NewC(code, format, arg)

	// Assertions
	if err == nil {
		t.Errorf("Expected error to be non-nil")
		return
	}

	// Type assertion
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("Expected error to be of type *Error, got %T", err)
		return
	}

	// Check if the code is set correctly
	if e.Code != code {
		t.Errorf("Expected error code to be %v, got %v", code, e.Code)
	}

	// Check if the error message is formatted correctly
	expectedMsg := fmt.Sprintf("[%s](%s): %s", code.Str(), code.Msg(), fmt.Sprintf(format, arg))
	if e.Error() != expectedMsg {
		t.Errorf("Expected error message to be '%s', got '%s'", expectedMsg, e.Error())
	}
}

func TestWrapSupportsErrorsIs(t *testing.T) {
	sentinel := stderrors.New("sentinel")

	err := Wrap(sentinel)

	if !stderrors.Is(err, sentinel) {
		t.Fatalf("Expected errors.Is to match wrapped sentinel")
	}
}

func TestWrapSupportsErrorsAs(t *testing.T) {
	err := WithC(stderrors.New("sentinel"), codes.AppValidation, "friendly message")

	var wrapped *Error
	if !stderrors.As(err, &wrapped) {
		t.Fatalf("Expected errors.As to find *Error")
	}

	if wrapped.Code != codes.AppValidation {
		t.Fatalf("Expected code %v, got %v", codes.AppValidation, wrapped.Code)
	}
}

func TestWithSupportsErrorsIsAcrossMultipleWraps(t *testing.T) {
	sentinel := stderrors.New("sentinel")

	err := With(WrapC(sentinel, codes.AppValidation), "handling user %s", "alice")

	if !stderrors.Is(err, sentinel) {
		t.Fatalf("Expected errors.Is to match sentinel through nested wrappers")
	}
}

func TestWithPreservesExistingCode(t *testing.T) {
	base := NewC(codes.AppValidation, "invalid value")

	err := With(base, "request failed")

	e, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected error to be of type *Error, got %T", err)
	}

	if e.Code != codes.AppValidation {
		t.Fatalf("Expected code %v, got %v", codes.AppValidation, e.Code)
	}
}

func TestWithCSetsProvidedCode(t *testing.T) {
	base := NewC(codes.AppValidation, "invalid value")

	err := WithC(base, codes.SrvInternal, "request failed")

	e, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected error to be of type *Error, got %T", err)
	}

	if e.Code != codes.SrvInternal {
		t.Fatalf("Expected code %v, got %v", codes.SrvInternal, e.Code)
	}
}

func TestWrapCRewrapsExistingErrorAndPreservesMessage(t *testing.T) {
	sentinel := stderrors.New("sentinel")
	base := With(sentinel, "first context")

	err := WrapC(base, codes.SrvInternal)

	if !stderrors.Is(err, sentinel) {
		t.Fatalf("Expected errors.Is to match original sentinel")
	}

	e, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected error to be of type *Error, got %T", err)
	}

	if e.Message != "first context" {
		t.Fatalf("Expected message to be preserved, got %q", e.Message)
	}

	if got := e.Error(); strings.Count(got, "first context") != 1 {
		t.Fatalf("Expected message to appear once, got %q", got)
	}

	if e.Code != codes.SrvInternal {
		t.Fatalf("Expected code %v, got %v", codes.SrvInternal, e.Code)
	}
}

func TestNilInputsReturnNil(t *testing.T) {
	if Wrap(nil) != nil {
		t.Fatal("Expected Wrap(nil) to return nil")
	}

	if WrapC(nil, codes.Unknown) != nil {
		t.Fatal("Expected WrapC(nil) to return nil")
	}

	if With(nil, "message") != nil {
		t.Fatal("Expected With(nil) to return nil")
	}

	if WithC(nil, codes.Unknown, "message") != nil {
		t.Fatal("Expected WithC(nil) to return nil")
	}
}

func TestUnwrapReturnsCause(t *testing.T) {
	cause := stderrors.New("root cause")
	err := Wrap(cause)

	unwrapped := stderrors.Unwrap(err)
	if unwrapped != cause {
		t.Fatalf("Expected unwrap to return cause %v, got %v", cause, unwrapped)
	}
}

func TestStackTraceHandlesMissingWrapper(t *testing.T) {
	err := &Error{Message: "missing wrapper"}

	if got := err.StackTraceStr(); got != "\n\n" {
		t.Fatalf("Expected empty stack trace string, got %q", got)
	}

	if got := err.StackTrace(); got.Kind() != slog.KindGroup {
		t.Fatalf("Expected empty slog group, got %v", got)
	}
}

func TestStackTraceIncludesMessageAndCause(t *testing.T) {
	err := With(stderrors.New("root cause"), "friendly message")

	got := err.(*Error).StackTraceStr()
	if !strings.Contains(got, "friendly message") {
		t.Fatalf("Expected stack trace to include message, got %q", got)
	}

	if !strings.Contains(got, "root cause") {
		t.Fatalf("Expected stack trace to include cause, got %q", got)
	}
}

func TestCodeOfReturnsCodeFromWrappedError(t *testing.T) {
	err := WithC(stderrors.New("root cause"), codes.AppValidation, "friendly message")

	if got := CodeOf(err); got != codes.AppValidation {
		t.Fatalf("Expected code %v, got %v", codes.AppValidation, got)
	}
}

func TestCodeOfReturnsNilForPlainError(t *testing.T) {
	if got := CodeOf(stderrors.New("plain")); got != nil {
		t.Fatalf("Expected nil code, got %v", got)
	}
}

func TestCodeOfSupportsCustomCodes(t *testing.T) {
	custom := codes.New("x001", "Custom failure.")
	err := NewC(custom, "custom message")

	if got := CodeOf(err); got != custom {
		t.Fatalf("Expected custom code %v, got %v", custom, got)
	}
}
