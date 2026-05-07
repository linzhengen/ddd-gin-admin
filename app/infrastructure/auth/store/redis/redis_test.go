package redis

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	addr = "127.0.0.1:6379"
)

func TestStore(t *testing.T) {
	// Skip if Redis is not available (e.g., CI, local dev without Redis)
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Skipf("Redis not available at %s: %v", addr, err)
	}
	_ = conn.Close()

	store := NewStore(&Config{
		Addr:      addr,
		DB:        1,
		KeyPrefix: "prefix",
	})

	defer func() { _ = store.Close() }()

	key := "test"
	ctx := context.Background()
	err = store.Set(ctx, key, 0)
	assert.Nil(t, err)

	b, err := store.Check(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	b, err = store.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)
}
