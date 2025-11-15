package errs

// HTTPError represents a domain error that carries an HTTP status code suitable for API responses.
type HTTPError struct {
	msg  string
	code int
}

// NewHTTPError constructs an HTTPError with the given message and status code.
func NewHTTPError(msg string, code int) *HTTPError {
	return &HTTPError{
		msg:  msg,
		code: code,
	}
}

// Error is a method that provides public behavior for the corresponding type.
func (e HTTPError) Error() string {
	return e.msg
}

// Code is a method that provides public behavior for the corresponding type.
func (e HTTPError) Code() int {
	return e.code
}

// DuplicatedURLError reports an attempt to shorten a URL that already exists.
// Handlers map this error to HTTP 409 Conflict.
type DuplicatedURLError struct {
	url string
}

// NewDuplicatedURLError creates a DuplicatedURLError for the given URL.
func NewDuplicatedURLError(url string) *DuplicatedURLError {
	return &DuplicatedURLError{
		url: url,
	}
}

// Error is a method that provides public behavior for the corresponding type.
func (e DuplicatedURLError) Error() string {
	return "duplicated URL: " + e.url
}
