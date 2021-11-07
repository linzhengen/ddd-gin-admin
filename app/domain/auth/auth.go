package auth

import (
	"errors"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Auth struct {
	AccessToken string
	TokenType   string
	ExpiresAt   int64
}

type RootUser struct {
	UserName string
	Password string
}
