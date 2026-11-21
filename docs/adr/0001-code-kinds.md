# 0001 — Codes declare an optional kind; transports map it to a status

- Status: accepted (2026-11-20)
- Issue: ERR-03 (#6), with API-01 in `jgolang/api`

## Context

Consumers kept per-package tables from error to HTTP status, and every new code had to be added to each table that could see it. Two proposals existed: a kind on the code (ERR-03) or an HTTP status on the code (API-01).

## Decision

- `codes.Code` can declare a `codes.Kind` through `WithKind`. The closed list approved by the owner is `validation`, `unauthorized`, `forbidden`, `not_found`, `conflict`, `infrastructure`, and `too_large` (added in v0.9.0, owner-approved 2026-11-21, jgolang/api#85).
- Kinds are read through the optional `codes.Kinded` interface, so `codes.Coder` and existing implementations do not change.
- `errors.KindOf(err)` returns the kind of `errors.CodeOf(err)`, never of another code in the chain, so a status derived from it always matches the reported code.
- This module stays free of HTTP. `jgolang/api` maps kinds to statuses (default 400, 401, 403, 404, 409, 500, and 413 for `too_large` from the api release that consumes v0.9.0; none or unknown is 500).
- Built-in codes carry kinds; `Unknown`, `Logic`, `Unhandled`, `Unexpected` and `AppProcessing` have none.

## Consequences

- One classification serves HTTP and any other transport.
- `codes.Code` gains a field, so a rebuilt `codes.New(...)` value no longer compares equal (`==`) to a built-in code that now has a kind. Noted in the v0.3.0 release notes.
- Adding a kind later is additive; removing or renaming one is a breaking change.
