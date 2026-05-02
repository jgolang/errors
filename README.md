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

## Stack Traces

`StackTraceStr` returns a text stack trace. `StackTrace` returns a `slog.Value` group that can be attached to structured logs. Stack traces can contain source paths and function names, so keep them in internal logs only.

```go
slog.Error("operation failed", "error", err, "stack", err.(*jerrors.Error).StackTrace())
```
