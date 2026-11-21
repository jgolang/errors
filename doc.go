// Package errors is an application error convention for Go services.
//
// An [Error] carries a stable code from package
// [github.com/jgolang/errors/codes], a user-safe message, and a stack trace
// for internal logs, while keeping the standard library errors.Is and errors.As
// chaining semantics. Create errors with [New], [NewC] or [Define], add
// context with [With], [WithC], [Wrap] or [WrapC], and read them back with
// [CodeOf], [KindOf] and [PublicMessage].
//
// Return [PublicMessage] to users or API clients; stack traces and wrapped
// causes can contain implementation details, so keep them in internal logs.
package errors
