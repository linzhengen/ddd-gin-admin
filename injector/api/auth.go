package api

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/linzhengen/ddd-gin-admin/app/domain/auth"
	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	authinfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/auth"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/auth/store/buntdb"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/auth/store/redis"
	"github.com/linzhengen/ddd-gin-admin/configs"
)

func InitAuth() (auth.Repository, func(), error) {
	cfg := configs.C.JWTAuth

	var opts []authinfra.Option
	opts = append(opts, authinfra.SetExpired(cfg.Expired))
	opts = append(opts, authinfra.SetSigningKey([]byte(cfg.SigningKey)))
	opts = append(opts, authinfra.SetKeyFunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
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
	opts = append(opts, authinfra.SetSigningMethod(method))
	opts = append(opts, authinfra.SetRootUser(configs.C.Root.UserName, configs.C.Root.Password))

	var store authinfra.Store
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

	auth := authinfra.NewRepository(store, opts...)
	cleanFunc := func() {
		//nolint:errcheck
		auth.Release()
	}
	return auth, cleanFunc, nil
}
