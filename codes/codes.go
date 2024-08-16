package codes

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
	Unknown   = Code{Code: "0001", Message: "Something went wrong. Please try again later."}
	Logic     = Code{Code: "0002", Message: "Logic error occurred."}
	Unhandled = Code{Code: "0003", Message: "Unexpected issue. Please contact support."}

	// Database errors
	ConnFailed = Code{Code: "d001", Message: "Database connection failed."}
	DBQuery    = Code{Code: "d002", Message: "Database query error."}

	// Network errors
	NetTimeout     = Code{Code: "n001", Message: "Network timeout."}
	NetUnavailable = Code{Code: "n002", Message: "Network unavailable."}

	// Server errors
	SrvInternal    = Code{Code: "s001", Message: "Internal server error."}
	SrvUnavailable = Code{Code: "s002", Message: "Server unavailable."}

	// Application Code errors
	AppValidation = Code{Code: "a001", Message: "Validation error."}
	AppProcessing = Code{Code: "a002", Message: "Processing error."}

	// Unexpected errors
	Unexpected = Code{Code: "u001", Message: "Unexpected error. Please contact support."}

	// Authentication errors
	AuthFailed       = Code{Code: "a003", Message: "Authentication failed."}
	AuthTokenExpired = Code{Code: "a004", Message: "Token expired."}

	// Authorization errors
	AuthForbidden = Code{Code: "f001", Message: "Access forbidden."}
	Unauthorized  = Code{Code: "f002", Message: "Unauthorized access."}

	// Validation errors
	ValRequired = Code{Code: "b001", Message: "Field missing."}
	ValFormat   = Code{Code: "b002", Message: "Invalid format."}

	// Configuration errors
	CfgInvalid = Code{Code: "c001", Message: "Invalid configuration."}
	CfgMissing = Code{Code: "c002", Message: "Configuration missing."}

	// Integration errors
	ServiceDown = Code{Code: "e001", Message: "Service down."}
	IntTimeout  = Code{Code: "e002", Message: "Integration timeout."}
)
