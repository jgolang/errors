package errors

// Code interface defines methods for retrieving error Code and message.
type Code interface {
	Str() string
	Msg() string
}

// BasicCode is a concrete implementation of the Code interface.
type BasicCode struct {
	Code    string
	Message string
}

// Code returns the Code as a string.
func (c BasicCode) Str() string {
	return c.Code
}

// Msg returns the message associated with the Code.
func (c BasicCode) Msg() string {
	return c.Message
}

// Define error Codes in the format A000
var (
	// General Code Errors
	ErrGenUnknown   = BasicCode{Code: "0001", Message: "Something went wrong. Please try again later."}
	ErrGenLogic     = BasicCode{Code: "0002", Message: "Logic error occurred."}
	ErrGenUnhandled = BasicCode{Code: "0003", Message: "Unexpected issue. Please contact support."}

	// Database errors
	ErrDBConnFailed = BasicCode{Code: "d001", Message: "Database connection failed."}
	ErrDBQuery      = BasicCode{Code: "d002", Message: "Database query error."}

	// Network errors
	ErrNetTimeout     = BasicCode{Code: "n001", Message: "Network timeout."}
	ErrNetUnavailable = BasicCode{Code: "n002", Message: "Network unavailable."}

	// Server errors
	ErrSrvInternal    = BasicCode{Code: "s001", Message: "Internal server error."}
	ErrSrvUnavailable = BasicCode{Code: "s002", Message: "Server unavailable."}

	// Application Code errors
	ErrAppValidation = BasicCode{Code: "a001", Message: "Validation error."}
	ErrAppProcessing = BasicCode{Code: "a002", Message: "Processing error."}

	// Unexpected errors
	ErrUnexpected = BasicCode{Code: "u001", Message: "Unexpected error. Please contact support."}

	// Authentication errors
	ErrAuthFailed       = BasicCode{Code: "a003", Message: "Authentication failed."}
	ErrAuthTokenExpired = BasicCode{Code: "a004", Message: "Token expired."}

	// Authorization errors
	ErrAuthForbidden    = BasicCode{Code: "f001", Message: "Access forbidden."}
	ErrAuthUnauthorized = BasicCode{Code: "f002", Message: "Unauthorized access."}

	// Validation errors
	ErrValRequired = BasicCode{Code: "b001", Message: "Field missing."}
	ErrValFormat   = BasicCode{Code: "b002", Message: "Invalid format."}

	// Configuration errors
	ErrCfgInvalid = BasicCode{Code: "c001", Message: "Invalid configuration."}
	ErrCfgMissing = BasicCode{Code: "c002", Message: "Configuration missing."}

	// Integration errors
	ErrIntServiceDown = BasicCode{Code: "e001", Message: "Service down."}
	ErrIntTimeout     = BasicCode{Code: "e002", Message: "Integration timeout."}
)
