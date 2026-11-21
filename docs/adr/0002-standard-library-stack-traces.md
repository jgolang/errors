# 0002 — Stack traces use the standard library

- Status: accepted (2026-11-21)
- Issue: ERR-07

## Context

`errors.Error` exported `Wrapper *goerrors.Error`, which tied the public API of this leaf module to `github.com/go-errors/errors`, its only third-party dependency. The module used it only to capture program counters and turn them into frames, which `runtime.Callers` and `runtime.CallersFrames` already do (KISS-01). `Wrapper` was deprecated in v0.7.0 in favor of `Detail`, `StackTrace` and `StackTraceStr`, and `jgolang/log` stopped reading it in v1.7.3.

## Decision

- Remove `Wrapper` and the go-errors dependency in v0.8.0. The module has no third-party dependency.
- `*Error` keeps the wrapped error and the captured program counters in unexported fields. An error has a stack trace exactly when it had a non-nil `Wrapper` before: errors built by `New`, `NewC`, `Wrap`, `With`, `WithC`, `WrapC` and `Wrap` of a `Define` sentinel; not `Define` sentinels or hand-built `&Error{...}` values.
- Constructors capture the stack with `runtime.Callers`, starting at their caller and keeping at most 50 frames (an unexported constant, as in go-errors).
- `StackTrace` and `StackTraceStr` resolve frames with `runtime.CallersFrames` on every call and cache nothing on the shared value, so concurrent reads are race-free (ERR-06). Frame names are split into package and name as go-errors did, so the output format is unchanged.

## Consequences

- Breaking change at v0: released as v0.8.0, whose notes name `jgolang/log` v1.7.3 as the minimum `log` version. Code that read `Wrapper` must use `Detail`, `StackTrace` or `StackTraceStr`.
- For frames the compiler did not inline, the output is the same as before. A function inlined into its caller now appears with its real name and line; go-errors repeated the callee or named another function (such as the standard library's `errors.New`) for such frames.
- Wrapping a `*goerrors.Error` no longer reuses its stack: every constructor captures a fresh stack at its call site.
- Creating an error saves one allocation (the `*goerrors.Error`); reading a stack trace allocates as many times as before.
