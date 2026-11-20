package codes

// Kind classifies a code by the kind of failure it reports. Transport layers
// such as jgolang/api map a kind to their own status; this package does not.
type Kind string

// Kinds a code can declare. A code with no kind reports an empty Kind.
const (
	KindValidation     Kind = "validation"
	KindUnauthorized   Kind = "unauthorized"
	KindForbidden      Kind = "forbidden"
	KindNotFound       Kind = "not_found"
	KindConflict       Kind = "conflict"
	KindInfrastructure Kind = "infrastructure"
)

// Kinded is implemented by a Coder that declares a Kind. It is optional, so
// existing Coder implementations keep compiling.
type Kinded interface {
	Kind() Kind
}

// WithKind returns a copy of c that declares kind.
func (c Code) WithKind(kind Kind) Code {
	c.kind = kind
	return c
}

// Kind returns the kind c declares, or an empty Kind.
func (c Code) Kind() Kind {
	return c.kind
}

// Code declares its kind through Kinded.
var _ Kinded = Code{}
