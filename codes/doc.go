// Package codes defines the stable error codes used by
// [github.com/jgolang/errors].
//
// A [Coder] pairs a machine-readable code string with a user-facing message,
// and may declare a [Kind] so a transport layer can map it to its own status
// without a per-code table. [New] creates application-specific codes; the
// built-in codes are package variables, to be treated as read-only.
package codes
