# errors

Application error convention for Go services.

This package keeps Go's standard error chaining semantics while adding stable error codes, user-safe public messages, and stack traces for internal logs. Use it when you want errors to carry both machine-readable context for the application and safe messages for users or API clients.

## Install

```sh
go get github.com/jgolang/errors
```

## Create Errors

```go
err := jerrors.NewC(codes.AppValidation, "invalid field %q", "email")
```

The error string includes the code, code message, and original error message:

```text
[a001](Validation error.): invalid field "email"
```

Text that repeats the code message is printed once: `NewC(codes.New("PLN-111", "task not found"), "task not found")` reads `[PLN-111](task not found)`. A wrapped error with the same code does not repeat its prefix either: `With(err, "loading board")` on the error above reads `[PLN-111](task not found) loading board`. This applies when `With`, `WithC` or `WrapC` wraps the coded error directly; behind another wrapper such as `fmt.Errorf` the inner text is printed as it is.

## Wrap Errors

```go
var ErrUserNotFound = stderrors.New("user not found")

err := jerrors.WithC(ErrUserNotFound, codes.AppValidation, "cannot load user %s", userID)
```

Use `With` when you want to add context and keep the existing code. `With` and `Wrap` find a code anywhere in the chain, including behind `fmt.Errorf("...: %w", err)`; without one they use `codes.Unknown`:

```go
err = jerrors.With(err, "request %s failed", requestID)
```

Use `WrapC` when you want to assign or replace a code without adding a new friendly message:

```go
err = jerrors.WrapC(err, codes.SrvInternal)
```

## errors.Is and errors.As

Errors created by this package implement `Unwrap`, so standard Go helpers work:

```go
if stderrors.Is(err, ErrUserNotFound) {
	// handle known cause
}

var wrapped *jerrors.Error
if stderrors.As(err, &wrapped) {
	slog.Info("request failed", "code", wrapped.Code.Str(), "stack", wrapped.StackTrace())
}
```

If this package is imported as `errors`, use its `Is` helper directly:

```go
if errors.Is(err, ErrUserNotFound) {
	// handle known cause
}
```

## Sentinels

`NewC` captures a stack trace when it runs, so a package-level variable built with it points at program start. Declare sentinels with `Define`, which captures no stack, and return them through `Wrap`, which captures the stack where they are returned:

```go
var ErrTitleRequired = jerrors.Define(CodeTitleRequired, "title is required")

func validate(title string) error {
	if title == "" {
		return jerrors.Wrap(ErrTitleRequired) // the stack trace starts here
	}
	return nil
}

stderrors.Is(validate(""), ErrTitleRequired) // true
jerrors.CodeOf(validate(""))                 // CodeTitleRequired
```

`Wrap` returns any other top-level `*jerrors.Error` unchanged, so wrapping twice does not add a second stack trace. Stack traces are captured with the standard `runtime` package and keep at most 50 frames; a function the Go compiler inlined still appears as its own frame, with its real name and line.

## Extract Codes

Use `CodeOf` to read the first non-nil code from an error chain without a type assertion. It looks behind other wrappers such as `fmt.Errorf` and skips errors whose code is nil:

```go
if code := jerrors.CodeOf(err); code != nil {
	fmt.Println(code.Str(), code.Msg())
}
```

## Public Messages

Use `PublicMessage` when returning an error to users or API clients. It returns the code message and avoids exposing internal causes, stack traces, file paths, SQL, credentials, or other implementation details.

```go
http.Error(w, jerrors.PublicMessage(err), http.StatusInternalServerError)
```

## Custom Codes

Built-in codes live in the `codes` package:

```go
err := jerrors.NewC(codes.NetTimeout, "upstream request timed out")
```

For application-specific codes, use `codes.New`:

```go
PaymentDeclined := codes.New("p001", "Payment declined.")
err := jerrors.NewC(PaymentDeclined, "processor rejected payment")
```

The built-in codes are package variables because Go cannot define const struct values. Treat them as read-only values.

## Code Uniqueness

Two failures that share a code string cannot be told apart by clients. Check your own codes in a test (here `CodeTitleRequired` and the others are codes your package declares with `codes.New`):

```go
func TestCodesAreUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, code := range []codes.Coder{CodeTitleRequired, CodeTaskNotFound, CodeProjectNotFound} {
		if seen[code.Str()] {
			t.Fatalf("code %s is declared twice", code.Str())
		}
		seen[code.Str()] = true
	}
}
```

`Example_codeUniqueness` in `example_uniqueness_test.go` shows the same check as a function that returns every duplicated code string.

## Kinds

A code can declare the kind of failure it reports, so a transport layer can map it to a status without a per-code table. This package does not know about HTTP; `jgolang/api` does the mapping.

```go
TitleRequired := codes.New("PLN-001", "title is required").WithKind(codes.KindValidation)

err := jerrors.NewC(TitleRequired, "title is empty")
jerrors.KindOf(fmt.Errorf("saving: %w", err)) // codes.KindValidation
```

Kinds: `KindValidation`, `KindUnauthorized`, `KindForbidden`, `KindNotFound`, `KindConflict`, `KindInfrastructure`. `KindOf` returns the kind of the code `CodeOf` returns, or an empty kind when that code declares none. A custom `Coder` declares a kind by implementing `codes.Kinded`.

Built-in codes carry kinds: `ValRequired`, `ValFormat` and `AppValidation` are validation; `AuthFailed`, `AuthTokenExpired` and `Unauthorized` are unauthorized; `AuthForbidden` is forbidden; the database, network, server, configuration and integration codes are infrastructure. `Unknown`, `Logic`, `Unhandled`, `Unexpected` and `AppProcessing` have no kind.

## Stack Traces

`StackTraceStr` returns a text stack trace. `StackTrace` returns a `slog.Value` group that can be attached to structured logs. `Detail` returns the text of the wrapped cause without the error's own code and friendly message (a wrapped `*Error` keeps its own `[code](message)` prefix). Stack traces and details can contain source paths, function names and other implementation details, so keep them in internal logs only.

The exported `Wrapper` field, deprecated in v0.7.0, was removed in v0.8.0 together with the `github.com/go-errors/errors` dependency: use `Detail` for the cause's text and `StackTrace` or `StackTraceStr` for the stack. Wrapping a `*goerrors.Error` now captures a fresh stack at the call site instead of reusing its stack. `jgolang/log` v1.7.3 or later is required with v0.8.0.

```go
slog.Error("operation failed", "error", err, "stack", err.(*jerrors.Error).StackTrace())
```

## License

Released under the [MIT License](LICENSE).
