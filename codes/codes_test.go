package codes

import "testing"

func TestCodeAccessors(t *testing.T) {
	code := New("x001", "test message")

	if got := code.Str(); got != "x001" {
		t.Fatalf("Expected code string %q, got %q", "x001", got)
	}

	if got := code.Msg(); got != "test message" {
		t.Fatalf("Expected code message %q, got %q", "test message", got)
	}
}
