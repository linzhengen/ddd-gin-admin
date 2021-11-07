package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"

	"github.com/linzhengen/ddd-gin-admin/app/domain/auth"
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

func SetKeyFunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyFunc = keyFunc
	}
}

func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

func SetRootUser(id, password string) Option {
	return func(o *options) {
		o.rootUser = auth.RootUser{
			UserName: id,
			Password: hash.MD5String(password),
		}
	}
}

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyFunc       jwt.Keyfunc
	expired       int
	tokenType     string
	rootUser      auth.RootUser
}

func NewRepository(store Store, opts ...Option) auth.Repository {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &repositoryImpl{
		opts:  &o,
		store: store,
	}
}

type repositoryImpl struct {
	opts  *options
	store Store
}

func (a *repositoryImpl) FindRootUser(ctx context.Context, userName string) *auth.RootUser {
	if userName == a.opts.rootUser.UserName {
		return &auth.RootUser{
			UserName: userName,
			Password: a.opts.rootUser.Password,
		}
	}
	return nil
}

func (a *repositoryImpl) GenerateToken(ctx context.Context, userID string) (*auth.Auth, error) {
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

	tokenInfo := &auth.Auth{
		ExpiresAt:   expiresAt,
		TokenType:   a.opts.tokenType,
		AccessToken: tokenString,
	}
	return tokenInfo, nil
}

func (a *repositoryImpl) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.opts.keyFunc)
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func (a *repositoryImpl) callStore(fn func(Store) error) error {
	if store := a.store; store != nil {
		return fn(store)
	}
	return nil
}

func (a *repositoryImpl) DestroyToken(ctx context.Context, tokenString string) error {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return err
	}

	// save black list token when use store
	return a.callStore(func(store Store) error {
		//nolint:gosimple
		expired := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
		return store.Set(ctx, tokenString, expired)
	})
}

func (a *repositoryImpl) ParseUserID(ctx context.Context, tokenString string) (string, error) {
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

func (a *repositoryImpl) Release() error {
	return a.callStore(func(store Store) error {
		return store.Close()
	})
}
