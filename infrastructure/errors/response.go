package errors

// ResponseError ...
type ResponseError struct {
	Code       int    // error code
	Message    string // error msg
	StatusCode int    // error status code
	ERR        error  // error
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return r.Message
}

// UnWrapResponse ...
func UnWrapResponse(err error) *ResponseError {
	if v, ok := err.(*ResponseError); ok {
		return v
	}
	return nil
}

// WrapResponse ...
func WrapResponse(err error, code int, msg string, status ...int) error {
	res := &ResponseError{
		Code:    code,
		Message: msg,
		ERR:     err,
	}
	if len(status) > 0 {
		res.StatusCode = status[0]
	}
	return res
}

// Wrap400Response ...
func Wrap400Response(err error, msg ...string) error {
	m := "400 error"
	if len(msg) > 0 {
		m = msg[0]
	}
	return WrapResponse(err, 400, m, 400)
}

// Wrap500Response ...
func Wrap500Response(err error, msg ...string) error {
	m := "500 error"
	if len(msg) > 0 {
		m = msg[0]
	}
	return WrapResponse(err, 500, m, 500)
}

// NewResponse ...
func NewResponse(code int, msg string, status ...int) error {
	res := &ResponseError{
		Code:    code,
		Message: msg,
	}
	if len(status) > 0 {
		res.StatusCode = status[0]
	}
	return res
}

// New400Response ...
func New400Response(msg string) error {
	return NewResponse(400, msg, 400)
}

// New500Response ...
func New500Response(msg string) error {
	return NewResponse(500, msg, 500)
}
