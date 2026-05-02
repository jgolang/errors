# errors

Small Go error helper for wrapping errors with stack traces, friendly messages, and application codes.

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

## Wrap Errors

```go
var ErrUserNotFound = stderrors.New("user not found")

err := jerrors.WithC(ErrUserNotFound, codes.AppValidation, "cannot load user %s", userID)
```

Use `With` when you want to add context and keep the existing code if the error was already created by this library:

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
if errors.Is(err, ErrUserNotFound) {
	// handle known cause
}

var wrapped *jerrors.Error
if errors.As(err, &wrapped) {
	slog.Info("request failed", "code", wrapped.Code.Str(), "stack", wrapped.StackTrace())
}
```

## Extract Codes

Use `CodeOf` to read the first code from an error chain without a type assertion:

```go
if code := jerrors.CodeOf(err); code != nil {
	fmt.Println(code.Str(), code.Msg())
}
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

## Stack Traces

`StackTraceStr` returns a text stack trace. `StackTrace` returns a `slog.Value` group that can be attached to structured logs.

```go
slog.Error("operation failed", "error", err, "stack", err.(*jerrors.Error).StackTrace())
```
