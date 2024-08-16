package errors

// Code interface defines methods for retrieving error code and message.
type Code interface {
	Str() string
	Msg() string
}

// BasicCode is a concrete implementation of the Code interface.
type BasicCode struct {
	code string
	msg  string
}

// Code returns the code as a string.
func (c BasicCode) Str() string {
	return c.code
}

// Msg returns the message associated with the code.
func (c BasicCode) Msg() string {
	return c.msg
}

// Define error codes in the format A000
var (
	// General Code Errors
	ErrGenUnknown   = BasicCode{code: "0001", msg: "Something went wrong. Please try again later."}
	ErrGenLogic     = BasicCode{code: "0002", msg: "Logic error occurred."}
	ErrGenUnhandled = BasicCode{code: "0003", msg: "Unexpected issue. Please contact support."}

	// Database errors
	ErrDBConnFailed = BasicCode{code: "d001", msg: "Database connection failed."}
	ErrDBQuery      = BasicCode{code: "d002", msg: "Database query error."}

	// Network errors
	ErrNetTimeout     = BasicCode{code: "n001", msg: "Network timeout."}
	ErrNetUnavailable = BasicCode{code: "n002", msg: "Network unavailable."}

	// Server errors
	ErrSrvInternal    = BasicCode{code: "s001", msg: "Internal server error."}
	ErrSrvUnavailable = BasicCode{code: "s002", msg: "Server unavailable."}

	// Application code errors
	ErrAppValidation = BasicCode{code: "a001", msg: "Validation error."}
	ErrAppProcessing = BasicCode{code: "a002", msg: "Processing error."}

	// Unexpected errors
	ErrUnexpected = BasicCode{code: "u001", msg: "Unexpected error. Please contact support."}

	// Authentication errors
	ErrAuthFailed       = BasicCode{code: "a003", msg: "Authentication failed."}
	ErrAuthTokenExpired = BasicCode{code: "a004", msg: "Token expired."}

	// Authorization errors
	ErrAuthForbidden    = BasicCode{code: "f001", msg: "Access forbidden."}
	ErrAuthUnauthorized = BasicCode{code: "f002", msg: "Unauthorized access."}

	// Validation errors
	ErrValRequired = BasicCode{code: "b001", msg: "Field missing."}
	ErrValFormat   = BasicCode{code: "b002", msg: "Invalid format."}

	// Configuration errors
	ErrCfgInvalid = BasicCode{code: "c001", msg: "Invalid configuration."}
	ErrCfgMissing = BasicCode{code: "c002", msg: "Configuration missing."}

	// Integration errors
	ErrIntServiceDown = BasicCode{code: "e001", msg: "Service down."}
	ErrIntTimeout     = BasicCode{code: "e002", msg: "Integration timeout."}
)
