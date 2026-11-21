package errors

import (
	stderrors "errors"
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"
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

func TestIsMatchesWrappedSentinel(t *testing.T) {
	sentinel := stderrors.New("sentinel")

	err := Wrap(sentinel)

	if !Is(err, sentinel) {
		t.Fatalf("Expected Is to match wrapped sentinel")
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

func TestPublicMessageReturnsCodeMessage(t *testing.T) {
	err := WithC(stderrors.New("database password leaked in cause"), codes.SrvInternal, "internal context")

	if got := PublicMessage(err); got != codes.SrvInternal.Msg() {
		t.Fatalf("Expected public message %q, got %q", codes.SrvInternal.Msg(), got)
	}
}

func TestPublicMessageFallsBackToUnknown(t *testing.T) {
	err := stderrors.New("internal detail")

	if got := PublicMessage(err); got != codes.Unknown.Msg() {
		t.Fatalf("Expected fallback message %q, got %q", codes.Unknown.Msg(), got)
	}
}

func TestWrapAndWithKeepCodeBehindFmtErrorf(t *testing.T) {
	notFound := codes.New("PLN-111", "task not found")
	sentinel := NewC(notFound, "task not found")

	tests := []struct {
		name string
		wrap func(error) error
	}{
		{name: "Wrap", wrap: Wrap},
		{name: "With", wrap: func(err error) error { return With(err, "context") }},
	}

	for _, tt := range tests {
		for _, depth := range []int{1, 2} {
			t.Run(fmt.Sprintf("%s/depth%d", tt.name, depth), func(t *testing.T) {
				var inner error = sentinel
				for i := 0; i < depth; i++ {
					inner = fmt.Errorf("level %d: %w", i, inner)
				}

				err := tt.wrap(inner)

				if got := CodeOf(err); got != notFound {
					t.Fatalf("Expected code %v, got %v", notFound, got)
				}
				if !stderrors.Is(err, sentinel) {
					t.Fatal("Expected errors.Is to match the sentinel")
				}
			})
		}
	}
}

func TestWrapCAssignsGivenCodeBehindFmtErrorf(t *testing.T) {
	base := NewC(codes.AppValidation, "invalid value")
	inner := fmt.Errorf("loading: %w", base)

	err := WrapC(inner, codes.SrvInternal)

	if got := CodeOf(err); got != codes.SrvInternal {
		t.Fatalf("Expected code %v, got %v", codes.SrvInternal, got)
	}
	if !stderrors.Is(err, base) {
		t.Fatal("Expected errors.Is to match the base error")
	}
}

func TestCodeOfSkipsNilCode(t *testing.T) {
	base := NewC(codes.AppValidation, "invalid value")
	err := WrapC(fmt.Errorf("context: %w", base), nil)

	if got := CodeOf(err); got != codes.AppValidation {
		t.Fatalf("Expected code %v, got %v", codes.AppValidation, got)
	}
}

func TestCodeOfReturnsNilWhenNoCodeInChain(t *testing.T) {
	err := WrapC(stderrors.New("plain"), nil)

	if got := CodeOf(err); got != nil {
		t.Fatalf("Expected nil code, got %v", got)
	}
}

func TestWrapAndWithUseUnknownWhenNoCodeInChain(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "Wrap over plain", err: Wrap(fmt.Errorf("x: %w", stderrors.New("plain")))},
		{name: "With over plain", err: With(fmt.Errorf("x: %w", stderrors.New("plain")), "context")},
		{name: "With over nil code", err: With(WrapC(stderrors.New("plain"), nil), "context")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CodeOf(tt.err); got != codes.Unknown {
				t.Fatalf("Expected code %v, got %v", codes.Unknown, got)
			}
		})
	}
}

func TestWrapAndWithDoNotPanicOnTypedNilError(t *testing.T) {
	inner := fmt.Errorf("x: %w", (*Error)(nil))

	if got := CodeOf(Wrap(inner)); got != codes.Unknown {
		t.Fatalf("Expected code %v from Wrap, got %v", codes.Unknown, got)
	}
	if got := CodeOf(With(inner, "context")); got != codes.Unknown {
		t.Fatalf("Expected code %v from With, got %v", codes.Unknown, got)
	}
}

var errDefinedTitleRequired = Define(codes.New("PLN-001", "title is required"), "title is required")

// returnDefinedSentinel is not inlined so it appears as its own stack frame.
//
//go:noinline
func returnDefinedSentinel() error {
	return Wrap(errDefinedTitleRequired)
}

func TestDefineHasNoStack(t *testing.T) {
	e, ok := errDefinedTitleRequired.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T", errDefinedTitleRequired)
	}
	if len(e.StackTrace().Group()) != 0 || e.StackTraceStr() != "\n\n" || e.Detail() != "" {
		t.Fatal("Expected a sentinel without a stack trace")
	}
	if e.Code.Str() != "PLN-001" || e.Message != "title is required" {
		t.Fatalf("Expected code and message to be kept, got %v %q", e.Code, e.Message)
	}
}

func TestWrapOfDefinedSentinel(t *testing.T) {
	err := returnDefinedSentinel()

	t.Run("stack starts at the caller", func(t *testing.T) {
		e, ok := err.(*Error)
		if !ok {
			t.Fatalf("Expected *Error, got %T", err)
		}
		frames := e.StackTrace().Group()
		if len(frames) == 0 {
			t.Fatal("Expected a stack trace")
		}
		if top := frames[0].Value.String(); !strings.HasSuffix(top, " (returnDefinedSentinel)") {
			t.Fatalf("Expected top frame returnDefinedSentinel, got %q", top)
		}
	})

	t.Run("errors.Is matches the sentinel", func(t *testing.T) {
		if !stderrors.Is(err, errDefinedTitleRequired) {
			t.Fatal("Expected errors.Is to match the sentinel")
		}
	})

	t.Run("errors.As exposes the sentinel code", func(t *testing.T) {
		var e *Error
		if !stderrors.As(err, &e) || e.Code.Str() != "PLN-001" {
			t.Fatalf("Expected errors.As to find code PLN-001, got %v", e)
		}
	})

	t.Run("CodeOf returns the sentinel code", func(t *testing.T) {
		if got := CodeOf(err); got.Str() != "PLN-001" {
			t.Fatalf("Expected code PLN-001, got %v", got)
		}
	})

	t.Run("Error reads like NewC", func(t *testing.T) {
		if got := err.Error(); got != "[PLN-001](title is required)" {
			t.Fatalf("Unexpected Error() text %q", got)
		}
	})

	t.Run("a second Wrap returns the same error", func(t *testing.T) {
		if Wrap(err) != err {
			t.Fatal("Expected Wrap of a traced sentinel to return it unchanged")
		}
	})
}

func TestWrapReturnsHandBuiltStacklessErrorUnchanged(t *testing.T) {
	handBuilt := &Error{Code: codes.AppValidation, Message: "invalid value"}

	if Wrap(handBuilt) != error(handBuilt) {
		t.Fatal("Expected Wrap to return a hand-built *Error unchanged")
	}
}

func TestWrapReturnsStackedErrorUnchanged(t *testing.T) {
	stacked := NewC(codes.AppValidation, "invalid value")

	if Wrap(stacked) != stacked {
		t.Fatal("Expected Wrap to return a stacked *Error unchanged")
	}
}

func TestKindOf(t *testing.T) {
	validation := codes.New("PLN-001", "title is required").WithKind(codes.KindValidation)

	tests := []struct {
		name string
		err  error
		want codes.Kind
	}{
		{"coded error", NewC(validation, "bad"), codes.KindValidation},
		{"behind fmt.Errorf", fmt.Errorf("x: %w", NewC(validation, "bad")), codes.KindValidation},
		{"too large behind fmt.Errorf", fmt.Errorf("x: %w", NewC(codes.New("UPL-001", "file is too large").WithKind(codes.KindTooLarge), "bad")), codes.KindTooLarge},
		{"code without kind", NewC(codes.New("x001", "m"), "bad"), ""},
		{"outer code without kind wins", WrapC(NewC(validation, "bad"), codes.New("x001", "m")), ""},
		{"custom Coder", NewC(plainCoder{}, "bad"), ""},
		{"custom Kinded Coder", NewC(kindedCoder{}, "bad"), codes.KindConflict},
		{"plain error", stderrors.New("plain"), ""},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KindOf(tt.err); got != tt.want {
				t.Fatalf("Expected kind %q, got %q", tt.want, got)
			}
		})
	}
}

type plainCoder struct{}

func (plainCoder) Str() string { return "c001" }
func (plainCoder) Msg() string { return "custom" }

type kindedCoder struct{ plainCoder }

func (kindedCoder) Kind() codes.Kind { return codes.KindConflict }

func TestErrorPrintsCodeTextOnce(t *testing.T) {
	notFound := codes.New("PLN-111", "task not found")

	tests := []struct {
		name string
		err  error
		want string
	}{
		{"NewC with the code text", NewC(notFound, "task not found"), "[PLN-111](task not found)"},
		{"NewC with other text", NewC(notFound, "task 42 is missing"), "[PLN-111](task not found): task 42 is missing"},
		{"Define sentinel", Define(notFound, "task not found"), "[PLN-111](task not found)"},
		{"traced Define sentinel", Wrap(Define(notFound, "task not found")), "[PLN-111](task not found)"},
		{"WithC drops a message equal to the code text", WithC(stderrors.New("io failure"), notFound, "task not found"), "[PLN-111](task not found): io failure"},
		{"WrapC drops a cause equal to the code text", WrapC(stderrors.New("task not found"), notFound), "[PLN-111](task not found)"},
		{"text differing in case prints twice", NewC(notFound, "Task not found"), "[PLN-111](task not found): Task not found"},
		{"code without message keeps the separator", NewC(codes.New("X", ""), ""), "[X]: "},
		{"nil code keeps the message", WithC(stderrors.New("io failure"), nil, "loading board"), " loading board: io failure"},
		{"With adds context", With(NewC(notFound, "task not found"), "loading board"), "[PLN-111](task not found) loading board"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("Expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestErrorPrintsChainCodePrefixOnce(t *testing.T) {
	notFound := codes.New("PLN-111", "task not found")
	conflict := codes.New("PLN-200", "task is locked")

	tests := []struct {
		name string
		err  error
		want string
	}{
		{"With over the code text", With(NewC(notFound, "task not found"), "loading board"), "[PLN-111](task not found) loading board"},
		{"With over a detail", With(NewC(notFound, "task 42 is missing"), "loading board"), "[PLN-111](task not found) loading board: task 42 is missing"},
		{"two Withs", With(With(NewC(notFound, "task 42 is missing"), "loading board"), "handler"), "[PLN-111](task not found) handler: loading board: task 42 is missing"},
		{"Wrap behind fmt.Errorf keeps the inner text", Wrap(fmt.Errorf("loading: %w", NewC(notFound, "task 42 is missing"))), "[PLN-111](task not found): loading: [PLN-111](task not found): task 42 is missing"},
		{"different codes keep both", WithC(NewC(notFound, "task 42 is missing"), conflict, "saving"), "[PLN-200](task is locked) saving: [PLN-111](task not found): task 42 is missing"},
		{"same code under a different layer is kept", WithC(WithC(NewC(notFound, "task 42 is missing"), conflict, "ctx"), notFound, "saving"), "[PLN-111](task not found) saving: [PLN-200](task is locked) ctx: [PLN-111](task not found): task 42 is missing"},
		{"empty code message deduped", With(NewC(codes.New("PLN-111", ""), "detail"), "ctx"), "[PLN-111] ctx: detail"},
		{"WrapC over NewC with the same code", WrapC(NewC(notFound, "task 42 is missing"), notFound), "[PLN-111](task not found): task 42 is missing"},
		{"WrapC over a With chain", WrapC(With(NewC(notFound, "task 42 is missing"), "ctx"), notFound), "[PLN-111](task not found) ctx: task 42 is missing"},
		{"WrapC of a Define sentinel with its own code", WrapC(Define(notFound, "task 42 is missing"), notFound), "[PLN-111](task not found) task 42 is missing: task 42 is missing"},
		{"WrapC over a traced sentinel", WrapC(Wrap(Define(notFound, "task 42 is missing")), notFound), "[PLN-111](task not found): task 42 is missing"},
		{"With repeating the inner message", With(With(NewC(notFound, "task 42"), "loading"), "loading"), "[PLN-111](task not found) loading: loading: task 42"},
		{"inner text with a leading space is kept", With(NewC(notFound, " padded"), "ctx"), "[PLN-111](task not found) ctx:  padded"},
		{"inner empty detail keeps its marker", With(NewC(notFound, ""), "ctx"), "[PLN-111](task not found) ctx: [PLN-111](task not found): "},
		{"nil code without message", WrapC(stderrors.New("plain"), nil), "plain"},
		{"nil code chain", With(WrapC(stderrors.New("plain"), nil), "context"), "[0001](Something went wrong. Please try again later.) context: plain"},
		{"empty wrapped text keeps the separator", WithC(stderrors.New(""), notFound, "ctx"), "[PLN-111](task not found) ctx: "},
		{"different message for the same code string", WithC(NewC(notFound, "detail"), codes.New("PLN-111", ""), "ctx"), "[PLN-111] ctx: [PLN-111](task not found): detail"},
		{"Unknown chain", With(New("io failure"), "ctx"), "[0001](Something went wrong. Please try again later.) ctx: io failure"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("Expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestStackTraceIsSafeForConcurrentReads(t *testing.T) {
	err := NewC(codes.AppValidation, "invalid value").(*Error)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_ = err.StackTrace()
			_ = err.StackTraceStr()
		}()
	}
	close(start)
	wg.Wait()
}

// newCAtKnownLine is not inlined so it appears as its own stack frame; it
// returns the file and line where NewC is called.
//
//go:noinline
func newCAtKnownLine() (err *Error, file string, line int) {
	_, file, line, _ = runtime.Caller(0)
	err = NewC(codes.AppValidation, "invalid value").(*Error)
	return err, file, line + 1
}

// withAtKnownLine is not inlined so it appears as its own stack frame; it
// returns the file and line where With is called.
//
//go:noinline
func withAtKnownLine() (err *Error, file string, line int) {
	_, file, line, _ = runtime.Caller(0)
	err = With(stderrors.New("root cause"), "friendly message").(*Error)
	return err, file, line + 1
}

func TestStackTraceTopFrameIsTheCaller(t *testing.T) {
	tests := []struct {
		name  string
		build func() (*Error, string, int)
		want  string
	}{
		{name: "NewC", build: newCAtKnownLine, want: "newCAtKnownLine"},
		{name: "With", build: withAtKnownLine, want: "withAtKnownLine"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err, file, line := test.build()
			location := file + ":" + strconv.Itoa(line)

			frames := err.StackTrace().Group()
			if len(frames) == 0 {
				t.Fatal("Expected a stack trace")
			}
			wantAttr := slog.String("frame_0", location+" ("+test.want+")")
			if !frames[0].Equal(wantAttr) {
				t.Fatalf("Expected top frame %q, got %q", wantAttr, frames[0])
			}

			wantStr := "\n\ngithub.com/jgolang/errors." + test.want + "\n·    " + location + "\n"
			if got := err.StackTraceStr(); !strings.Contains(got, wantStr) {
				t.Fatalf("Expected StackTraceStr to contain\n%s\ngot\n%s", wantStr, got)
			}
		})
	}
}

// newInlinable is small enough for the compiler to inline into its caller.
func newInlinable() error {
	return New("inlined")
}

// callNewInlinable is not inlined; it returns the line where it calls
// newInlinable.
//
//go:noinline
func callNewInlinable() (err *Error, line int) {
	_, _, line, _ = runtime.Caller(0)
	err = newInlinable().(*Error)
	return err, line + 1
}

func TestStackTraceNamesInlinedCallers(t *testing.T) {
	err, line := callNewInlinable()

	frames := err.StackTrace().Group()
	if len(frames) < 2 {
		t.Fatalf("Expected at least 2 frames, got %d", len(frames))
	}
	if got := frames[0].Value.String(); !strings.HasSuffix(got, " (newInlinable)") {
		t.Fatalf("Expected frame 0 newInlinable, got %q", got)
	}
	want := ":" + strconv.Itoa(line) + " (callNewInlinable)"
	if got := frames[1].Value.String(); !strings.HasSuffix(got, want) {
		t.Fatalf("Expected frame 1 to end with %q, got %q", want, got)
	}
}

// newAtDepth is not inlined so each recursive call is its own stack frame.
//
//go:noinline
func newAtDepth(depth int) error {
	if depth == 0 {
		return New("deep")
	}
	return newAtDepth(depth - 1)
}

func TestStackTraceKeepsAtMostFiftyFrames(t *testing.T) {
	err := newAtDepth(60).(*Error)

	if got := len(err.StackTrace().Group()); got != 50 {
		t.Fatalf("Expected 50 frames, got %d", got)
	}
	if got := strings.Count(err.StackTraceStr(), "\n·    "); got != 50 {
		t.Fatalf("Expected 50 frames in StackTraceStr, got %d", got)
	}
}

func TestNilErrorHasNoStackOrText(t *testing.T) {
	var err *Error

	if got := len(err.StackTrace().Group()); got != 0 {
		t.Fatalf("Expected an empty stack trace, got %d frames", got)
	}
	if got := err.StackTraceStr(); got != "\n\n" {
		t.Fatalf("Expected empty stack trace string, got %q", got)
	}
	if got := err.Error(); got != "" {
		t.Fatalf("Expected empty text, got %q", got)
	}
}

func TestDetailReturnsWrappedCauseText(t *testing.T) {
	inner := NewC(codes.AppValidation, "title is empty")
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "New", err: New("disk %s", "full"), want: "disk full"},
		{name: "NewC", err: NewC(codes.AppValidation, "title is empty"), want: "title is empty"},
		{name: "Wrap of a plain error", err: Wrap(stderrors.New("timeout")), want: "timeout"},
		{name: "With", err: With(stderrors.New("timeout"), "loading board"), want: "timeout"},
		{name: "WithC", err: WithC(stderrors.New("timeout"), codes.AppProcessing, "loading board"), want: "timeout"},
		{name: "WrapC of an *Error without cause", err: WrapC(inner, codes.AppProcessing), want: inner.Error()},
		{name: "WrapC of an *Error wraps its cause", err: WrapC(With(stderrors.New("timeout"), "loading"), codes.AppProcessing), want: "timeout"},
		{name: "With of an *Error", err: With(inner, "saving"), want: inner.Error()},
		{name: "Wrap of a Define sentinel", err: Wrap(errDefinedTitleRequired), want: "title is required"},
		{name: "Define sentinel has no detail", err: errDefinedTitleRequired, want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.err.(*Error).Detail()

			if got != test.want {
				t.Fatalf("Expected %q, got %q", test.want, got)
			}
		})
	}
}

func TestDetailOfNilErrorIsEmpty(t *testing.T) {
	var err *Error

	if got := err.Detail(); got != "" {
		t.Fatalf("Expected empty detail, got %q", got)
	}
}

// requireComparable fails to compile when T is not comparable.
func requireComparable[T comparable]() {}

func TestErrorIsComparable(t *testing.T) {
	requireComparable[Error]()

	err := NewC(codes.AppValidation, "invalid value").(*Error)
	if *err != *err {
		t.Fatal("Expected an Error value to equal itself")
	}
}
