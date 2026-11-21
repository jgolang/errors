package errors

import (
	stderrors "errors"
	"testing"

	"github.com/jgolang/errors/codes"
)

func BenchmarkNewC(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = NewC(codes.AppValidation, "invalid value")
	}
}

func BenchmarkWith(b *testing.B) {
	cause := stderrors.New("timeout")
	b.ReportAllocs()
	for b.Loop() {
		_ = With(cause, "loading board")
	}
}

func BenchmarkStackTrace(b *testing.B) {
	err := NewC(codes.AppValidation, "invalid value").(*Error)
	b.ReportAllocs()
	for b.Loop() {
		_ = err.StackTrace()
	}
}

func BenchmarkStackTraceStr(b *testing.B) {
	err := NewC(codes.AppValidation, "invalid value").(*Error)
	b.ReportAllocs()
	for b.Loop() {
		_ = err.StackTraceStr()
	}
}
