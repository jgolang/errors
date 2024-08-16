package errors

import (
	"fmt"
	"testing"
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
	if e.Code != ErrGenUnknown {
		t.Errorf("Expected error code to be ErrGenUnknown, got %v", e.Code)
	}

	// Check if the error message is formatted correctly
	expectedMsg := fmt.Sprintf("[%s](%s): %s", ErrGenUnknown.Str(), ErrGenUnknown.Msg(), fmt.Sprintf(format, arg))
	if e.Error() != expectedMsg {
		t.Errorf("Expected error message to be '%s', got '%s'", expectedMsg, e.Error())
	}
}

func TestNewError(t *testing.T) {
	// Test case parameters
	code := ErrAppValidation // Example code, can be replaced with any code
	format := "Validation failed for field: %s"
	arg := "username"

	// Create a new error using the NewError function
	err := NewError(code, format, arg)

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
