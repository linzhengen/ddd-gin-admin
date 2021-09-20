package errors

import "fmt"

type ResponseError struct {
	Code       int
	Message    string
	StatusCode int
	ERR        error
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return r.Message
}

func UnWrapResponse(err error) *ResponseError {
	if v, ok := err.(*ResponseError); ok {
		return v
	}
	return nil
}

func WrapResponse(err error, code, statusCode int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:       code,
		Message:    fmt.Sprintf(msg, args...),
		ERR:        err,
		StatusCode: statusCode,
	}
	return res
}

func Wrap400Response(err error, msg string, args ...interface{}) error {
	return WrapResponse(err, 400, 400, msg, args...)
}

func Wrap500Response(err error, msg string, args ...interface{}) error {
	return WrapResponse(err, 500, 500, msg, args...)
}

func NewResponse(code, statusCode int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:       code,
		Message:    fmt.Sprintf(msg, args...),
		StatusCode: statusCode,
	}
	return res
}

func New400Response(msg string, args ...interface{}) error {
	return NewResponse(400, 400, msg, args...)
}

func New500Response(msg string, args ...interface{}) error {
	return NewResponse(500, 500, msg, args...)
}
