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
	ErrGenUnknown   = BasicCode{code: "001", msg: "Something went wrong. Please try again later."}
	ErrGenLogic     = BasicCode{code: "002", msg: "Logic error occurred."}
	ErrGenUnhandled = BasicCode{code: "003", msg: "Unexpected issue. Please contact support."}

	// Database errors
	ErrDBConnFailed = BasicCode{code: "D001", msg: "Database connection failed."}
	ErrDBQuery      = BasicCode{code: "D002", msg: "Database query error."}

	// Network errors
	ErrNetTimeout     = BasicCode{code: "N001", msg: "Network timeout."}
	ErrNetUnavailable = BasicCode{code: "N002", msg: "Network unavailable."}

	// Server errors
	ErrSrvInternal    = BasicCode{code: "S001", msg: "Internal server error."}
	ErrSrvUnavailable = BasicCode{code: "S002", msg: "Server unavailable."}

	// Application code errors
	ErrAppValidation = BasicCode{code: "A001", msg: "Validation error."}
	ErrAppProcessing = BasicCode{code: "A002", msg: "Processing error."}

	// Unexpected errors
	ErrUnexpected = BasicCode{code: "U001", msg: "Unexpected error. Please contact support."}

	// Authentication errors
	ErrAuthFailed       = BasicCode{code: "A003", msg: "Authentication failed."}
	ErrAuthTokenExpired = BasicCode{code: "A004", msg: "Token expired."}

	// Authorization errors
	ErrAuthForbidden    = BasicCode{code: "O001", msg: "Access forbidden."}
	ErrAuthUnauthorized = BasicCode{code: "O002", msg: "Unauthorized access."}

	// Validation errors
	ErrValRequired = BasicCode{code: "V001", msg: "Field missing."}
	ErrValFormat   = BasicCode{code: "V002", msg: "Invalid format."}

	// Configuration errors
	ErrCfgInvalid = BasicCode{code: "C001", msg: "Invalid configuration."}
	ErrCfgMissing = BasicCode{code: "C002", msg: "Configuration missing."}

	// Integration errors
	ErrIntServiceDown = BasicCode{code: "I001", msg: "Service down."}
	ErrIntTimeout     = BasicCode{code: "I002", msg: "Integration timeout."}
)
