package errs

type HTTPError struct {
	msg  string
	code int
}

func NewHTTPError(msg string, code int) *HTTPError {
	return &HTTPError{
		msg:  msg,
		code: code,
	}
}

func (e HTTPError) Error() string {
	return e.msg
}

func (e HTTPError) Code() int {
	return e.code
}
