package api

import (
	"github.com/golang-jwt/jwt"
	auth2 "github.com/linzhengen/ddd-gin-admin/app/domain/auth"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/auth/store/buntdb"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/auth/store/redis"
	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/pkg/auth/jwtauth"
)

func InitAuth() (auth2.Author, func(), error) {
	cfg := configs.C.JWTAuth

	var opts []jwtauth.Option
	opts = append(opts, jwtauth.SetExpired(cfg.Expired))
	opts = append(opts, jwtauth.SetSigningKey([]byte(cfg.SigningKey)))
	opts = append(opts, jwtauth.SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth2.ErrInvalidToken
		}
		return []byte(cfg.SigningKey), nil
	}))

	var method jwt.SigningMethod
	switch cfg.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}
	opts = append(opts, jwtauth.SetSigningMethod(method))

	var store auth2.Store
	switch cfg.Store {
	case "redis":
		rcfg := configs.C.Redis
		store = redis.NewStore(&redis.Config{
			Addr:      rcfg.Addr,
			Password:  rcfg.Password,
			DB:        cfg.RedisDB,
			KeyPrefix: cfg.RedisPrefix,
		})
	default:
		s, err := buntdb.NewStore(cfg.FilePath)
		if err != nil {
			return nil, nil, err
		}
		store = s
	}

	auth := jwtauth.New(store, opts...)
	cleanFunc := func() {
		//nolint:errcheck
		auth.Release()
	}
	return auth, cleanFunc, nil
}
