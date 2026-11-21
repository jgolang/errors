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

func TestWithKindSetsKindAndKeepsCode(t *testing.T) {
	code := New("PLN-001", "title is required").WithKind(KindValidation)

	if code.Kind() != KindValidation {
		t.Fatalf("Expected kind %q, got %q", KindValidation, code.Kind())
	}
	if code.Str() != "PLN-001" || code.Msg() != "title is required" {
		t.Fatalf("Expected code and message to be kept, got %q %q", code.Str(), code.Msg())
	}
}

func TestKindTooLarge(t *testing.T) {
	code := New("UPL-001", "file is too large").WithKind(KindTooLarge)

	if got := code.Kind(); got != KindTooLarge {
		t.Fatalf("Expected kind %q, got %q", KindTooLarge, got)
	}
	if got := string(KindTooLarge); got != "too_large" {
		t.Fatalf("Expected kind string %q, got %q", "too_large", got)
	}
}

func TestNewCodeHasNoKind(t *testing.T) {
	if got := New("x001", "m").Kind(); got != "" {
		t.Fatalf("Expected no kind, got %q", got)
	}
}

func TestBuiltInCodeKinds(t *testing.T) {
	tests := []struct {
		name string
		code Code
		want Kind
	}{
		{"ValRequired", ValRequired, KindValidation},
		{"ValFormat", ValFormat, KindValidation},
		{"AppValidation", AppValidation, KindValidation},
		{"AuthFailed", AuthFailed, KindUnauthorized},
		{"AuthTokenExpired", AuthTokenExpired, KindUnauthorized},
		{"Unauthorized", Unauthorized, KindUnauthorized},
		{"AuthForbidden", AuthForbidden, KindForbidden},
		{"ConnFailed", ConnFailed, KindInfrastructure},
		{"DBQuery", DBQuery, KindInfrastructure},
		{"NetTimeout", NetTimeout, KindInfrastructure},
		{"NetUnavailable", NetUnavailable, KindInfrastructure},
		{"SrvInternal", SrvInternal, KindInfrastructure},
		{"SrvUnavailable", SrvUnavailable, KindInfrastructure},
		{"CfgInvalid", CfgInvalid, KindInfrastructure},
		{"CfgMissing", CfgMissing, KindInfrastructure},
		{"ServiceDown", ServiceDown, KindInfrastructure},
		{"IntTimeout", IntTimeout, KindInfrastructure},
		{"Unknown", Unknown, ""},
		{"Logic", Logic, ""},
		{"Unhandled", Unhandled, ""},
		{"Unexpected", Unexpected, ""},
		{"AppProcessing", AppProcessing, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Kind(); got != tt.want {
				t.Fatalf("Expected kind %q, got %q", tt.want, got)
			}
		})
	}
}
