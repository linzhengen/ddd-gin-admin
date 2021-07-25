package jwtauth

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/linzhengen/ddd-gin-admin/pkg/auth"
)

const defaultKey = "ddd-gin-admin"

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       7200,
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    []byte(defaultKey),
	keyFunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(defaultKey), nil
	},
}

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyFunc       jwt.Keyfunc
	expired       int
	tokenType     string
}

type Option func(*options)

func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

// SetKeyfunc 设定验证key的回调函数
func SetKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyFunc = keyFunc
	}
}

func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

func New(store Store, opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &JWTAuth{
		opts:  &o,
		store: store,
	}
}

type JWTAuth struct {
	opts  *options
	store Store
}

func (a *JWTAuth) GenerateToken(ctx context.Context, userID string) (auth.TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.opts.expired) * time.Second).Unix()

	token := jwt.NewWithClaims(a.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   userID,
	})

	tokenString, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return nil, err
	}

	tokenInfo := &tokenInfo{
		ExpiresAt:   expiresAt,
		TokenType:   a.opts.tokenType,
		AccessToken: tokenString,
	}
	return tokenInfo, nil
}

func (a *JWTAuth) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.opts.keyFunc)
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func (a *JWTAuth) callStore(fn func(Store) error) error {
	if store := a.store; store != nil {
		return fn(store)
	}
	return nil
}

func (a *JWTAuth) DestroyToken(ctx context.Context, tokenString string) error {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return err
	}

	// 如果设定了存储，则将未过期的令牌放入
	return a.callStore(func(store Store) error {
		expired := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
		return store.Set(ctx, tokenString, expired)
	})
}

func (a *JWTAuth) ParseUserID(ctx context.Context, tokenString string) (string, error) {
	if tokenString == "" {
		return "", auth.ErrInvalidToken
	}

	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", err
	}

	err = a.callStore(func(store Store) error {
		if exists, err := store.Check(ctx, tokenString); err != nil {
			return err
		} else if exists {
			return auth.ErrInvalidToken
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

func (a *JWTAuth) Release() error {
	return a.callStore(func(store Store) error {
		return store.Close()
	})
}
