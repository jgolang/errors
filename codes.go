package errors

// Code interface defines methods for retrieving error Code and message.
type Coder interface {
	Str() string
	Msg() string
}

// Code is a concrete implementation of the Code interface.
type Code struct {
	Code    string
	Message string
}

// Code returns the Code as a string.
func (c Code) Str() string {
	return c.Code
}

// Msg returns the message associated with the Code.
func (c Code) Msg() string {
	return c.Message
}

// Define error Codes in the format A000
var (
	// General Code Errors
	ErrGenUnknown   = Code{Code: "0001", Message: "Something went wrong. Please try again later."}
	ErrGenLogic     = Code{Code: "0002", Message: "Logic error occurred."}
	ErrGenUnhandled = Code{Code: "0003", Message: "Unexpected issue. Please contact support."}

	// Database errors
	ErrDBConnFailed = Code{Code: "d001", Message: "Database connection failed."}
	ErrDBQuery      = Code{Code: "d002", Message: "Database query error."}

	// Network errors
	ErrNetTimeout     = Code{Code: "n001", Message: "Network timeout."}
	ErrNetUnavailable = Code{Code: "n002", Message: "Network unavailable."}

	// Server errors
	ErrSrvInternal    = Code{Code: "s001", Message: "Internal server error."}
	ErrSrvUnavailable = Code{Code: "s002", Message: "Server unavailable."}

	// Application Code errors
	ErrAppValidation = Code{Code: "a001", Message: "Validation error."}
	ErrAppProcessing = Code{Code: "a002", Message: "Processing error."}

	// Unexpected errors
	ErrUnexpected = Code{Code: "u001", Message: "Unexpected error. Please contact support."}

	// Authentication errors
	ErrAuthFailed       = Code{Code: "a003", Message: "Authentication failed."}
	ErrAuthTokenExpired = Code{Code: "a004", Message: "Token expired."}

	// Authorization errors
	ErrAuthForbidden    = Code{Code: "f001", Message: "Access forbidden."}
	ErrAuthUnauthorized = Code{Code: "f002", Message: "Unauthorized access."}

	// Validation errors
	ErrValRequired = Code{Code: "b001", Message: "Field missing."}
	ErrValFormat   = Code{Code: "b002", Message: "Invalid format."}

	// Configuration errors
	ErrCfgInvalid = Code{Code: "c001", Message: "Invalid configuration."}
	ErrCfgMissing = Code{Code: "c002", Message: "Configuration missing."}

	// Integration errors
	ErrIntServiceDown = Code{Code: "e001", Message: "Service down."}
	ErrIntTimeout     = Code{Code: "e002", Message: "Integration timeout."}
)
