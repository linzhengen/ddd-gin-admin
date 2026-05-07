package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr      string
	DB        int
	Password  string
	KeyPrefix string
}

func NewStore(cfg *Config) *Store {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	return &Store{
		cli:    cli,
		prefix: cfg.KeyPrefix,
	}
}

func NewStoreWithClient(cli *redis.Client, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
	}
}

func NewStoreWithClusterClient(cli *redis.ClusterClient, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
	}
}

type redisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	TxPipeline() redis.Pipeliner
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
}

type Store struct {
	cli    redisClient
	prefix string
}

func (s *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", s.prefix, key)
}

func (s *Store) Set(ctx context.Context, tokenString string, expiration time.Duration) error {
	cmd := s.cli.Set(ctx, s.wrapperKey(tokenString), "1", expiration)
	return cmd.Err()
}

func (s *Store) Delete(ctx context.Context, tokenString string) (bool, error) {
	cmd := s.cli.Del(ctx, s.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) Check(ctx context.Context, tokenString string) (bool, error) {
	cmd := s.cli.Exists(ctx, s.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) Close() error {
	return s.cli.Close()
}
