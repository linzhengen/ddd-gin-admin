package auth

import (
	"errors"
)

// Deprecated: Use errors.ErrInvalidToken instead.
// ErrInvalidToken is kept for backward compatibility, delegates to the domain errors package.
var ErrInvalidToken = errors.New("invalid token")

type Auth struct {
	AccessToken string
	TokenType   string
	ExpiresAt   int64
}

type RootUser struct {
	UserName string
	Password string
}
