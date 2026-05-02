package codes

// Coder exposes a stable code string and a user-facing message.
type Coder interface {
	Str() string
	Msg() string
}

// Code is the default Coder implementation.
type Code struct {
	code    string
	message string
}

// New creates a Code with the supplied stable code and user-facing message.
func New(code, message string) Code {
	return Code{code: code, message: message}
}

// Str returns the Code as a string.
func (c Code) Str() string {
	return c.code
}

// Msg returns the message associated with the Code.
func (c Code) Msg() string {
	return c.message
}

// Built-in codes are variables because Go cannot define const struct values.
// Treat them as read-only package values; use New to create custom codes.
var (
	// General Code Errors
	Unknown   = New("0001", "Something went wrong. Please try again later.")
	Logic     = New("0002", "Logic error occurred.")
	Unhandled = New("0003", "Unexpected issue. Please contact support.")

	// Database errors
	ConnFailed = New("d001", "Database connection failed.")
	DBQuery    = New("d002", "Database query error.")

	// Network errors
	NetTimeout     = New("n001", "Network timeout.")
	NetUnavailable = New("n002", "Network unavailable.")

	// Server errors
	SrvInternal    = New("s001", "Internal server error.")
	SrvUnavailable = New("s002", "Server unavailable.")

	// Application Code errors
	AppValidation = New("a001", "Validation error.")
	AppProcessing = New("a002", "Processing error.")

	// Unexpected errors
	Unexpected = New("u001", "Unexpected error. Please contact support.")

	// Authentication errors
	AuthFailed       = New("a003", "Authentication failed.")
	AuthTokenExpired = New("a004", "Token expired.")

	// Authorization errors
	AuthForbidden = New("f001", "Access forbidden.")
	Unauthorized  = New("f002", "Unauthorized access.")

	// Validation errors
	ValRequired = New("b001", "Field missing.")
	ValFormat   = New("b002", "Invalid format.")

	// Configuration errors
	CfgInvalid = New("c001", "Invalid configuration.")
	CfgMissing = New("c002", "Configuration missing.")

	// Integration errors
	ServiceDown = New("e001", "Service down.")
	IntTimeout  = New("e002", "Integration timeout.")
)
