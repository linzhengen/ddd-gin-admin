package errors

import (
	"github.com/pkg/errors"
)

var (
	New          = errors.New
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

var (
	ErrBadRequest              = New400Response("ErrBadRequest")
	ErrInvalidParent           = New400Response("ErrInvalidParent")
	ErrNotAllowDeleteWithChild = New400Response("ErrNotAllowDeleteWithChild")
	ErrNotAllowDelete          = New400Response("ErrNotAllowDelete")
	ErrInvalidUserName         = New400Response("ErrInvalidUserName")
	ErrInvalidPassword         = New400Response("ErrInvalidPassword")
	ErrInvalidUser             = New400Response("ErrInvalidUser")
	ErrUserDisable             = New400Response("ErrUserDisable")

	ErrNoPerm          = NewResponse(401, 401, "ErrNoPerm")
	ErrInvalidToken    = NewResponse(9999, 401, "ErrInvalidToken")
	ErrNotFound        = NewResponse(404, 404, "ErrNotFound")
	ErrMethodNotAllow  = NewResponse(405, 405, "ErrMethodNotAllow")
	ErrTooManyRequests = NewResponse(429, 429, "ErrTooManyRequests")
	ErrInternalServer  = NewResponse(500, 500, "ErrInternalServer")
)
