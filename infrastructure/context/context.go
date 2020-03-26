package context

import (
	"context"
)

// context keys.
type (
	transCtx     struct{}
	transLockCtx struct{}
	userIDCtx    struct{}
	traceIDCtx   struct{}
)

// NewTrans ...
func NewTrans(ctx context.Context, trans interface{}) context.Context {
	return context.WithValue(ctx, transCtx{}, trans)
}

// FromTrans ...
func FromTrans(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(transCtx{})
	return v, v != nil
}

// NewTransLock ...
func NewTransLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, transLockCtx{}, struct{}{})
}

// FromTransLock ...
func FromTransLock(ctx context.Context) bool {
	v := ctx.Value(transLockCtx{})
	return v != nil
}

// NewUserID ...
func NewUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userID)
}

// FromUserID ...
func FromUserID(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s, s != ""
		}
	}
	return "", false
}

// NewTraceID ...
func NewTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDCtx{}, traceID)
}

// FromTraceID ...
func FromTraceID(ctx context.Context) (string, bool) {
	v := ctx.Value(traceIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s, s != ""
		}
	}
	return "", false
}
