package errors

import "github.com/pkg/errors"

// error func alias
var (
	New          = errors.New
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

// define errors.
var (
	ErrBadRequest              = New400Response("ErrBadRequest")
	ErrInvalidParent           = New400Response("ErrInvalidParent")
	ErrNotAllowDeleteWithChild = New400Response("ErrNotAllowDeleteWithChild")
	ErrNotAllowDelete          = New400Response("ErrNotAllowDelete")
	ErrInvalidUserName         = New400Response("ErrInvalidUserName")
	ErrInvalidPassword         = New400Response("ErrInvalidPassword")
	ErrInvalidUser             = New400Response("ErrInvalidUser")
	ErrUserDisable             = New400Response("ErrUserDisable")

	ErrNoPerm          = NewResponse(401, "ErrNoPerm", 401)
	ErrInvalidToken    = NewResponse(9999, "ErrInvalidToken", 401)
	ErrNotFound        = NewResponse(404, "ErrNotFound", 404)
	ErrMethodNotAllow  = NewResponse(405, "ErrMethodNotAllow", 405)
	ErrTooManyRequests = NewResponse(429, "ErrTooManyRequests", 429)
	ErrInternalServer  = NewResponse(500, "ErrInternalServer", 500)
)
