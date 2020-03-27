package logger

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

// keys.
const (
	TraceIDKey      = "trace_id"
	UserIDKey       = "user_id"
	SpanTitleKey    = "span_title"
	SpanFunctionKey = "span_function"
	VersionKey      = "version"
)

// TraceIDFunc ...
type TraceIDFunc func() string

var (
	version     string
	traceIDFunc TraceIDFunc
)

// Logger logrus Logger
type Logger = logrus.Logger

// StandardLogger standard logger.
func StandardLogger() *Logger {
	return logrus.StandardLogger()
}

// SetLevel log level.
func SetLevel(level int) {
	logrus.SetLevel(logrus.Level(level))
}

// SetFormatter logger formatter.
func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
}

// SetOutput set output.
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// SetVersion set version.
func SetVersion(v string) {
	version = v
}

// SetTraceIDFunc 设定追踪ID的处理函数
func SetTraceIDFunc(fn TraceIDFunc) {
	traceIDFunc = fn
}

func getTraceID() string {
	if traceIDFunc != nil {
		return traceIDFunc()
	}
	return ""
}

type (
	traceIDContextKey struct{}
	spanIDContextKey  struct{}
	userIDContextKey  struct{}
)

// NewTraceIDContext setup trace id into context.
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDContextKey{}, traceID)
}

// FromTraceIDContext get trace id from context.
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return getTraceID()
}

// NewUserIDContext create user id context.
func NewUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

// FromUserIDContext get user id from context.
func FromUserIDContext(ctx context.Context) string {
	v := ctx.Value(userIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type spanOptions struct {
	Title    string
	FuncName string
}

// SpanOption span options.
type SpanOption func(*spanOptions)

// SetSpanTitle set span title.
func SetSpanTitle(title string) SpanOption {
	return func(o *spanOptions) {
		o.Title = title
	}
}

// SetSpanFuncName set span func name.
func SetSpanFuncName(funcName string) SpanOption {
	return func(o *spanOptions) {
		o.FuncName = funcName
	}
}

// StartSpan start span.
func StartSpan(ctx context.Context, opts ...SpanOption) *Entry {
	if ctx == nil {
		ctx = context.Background()
	}

	var o spanOptions
	for _, opt := range opts {
		opt(&o)
	}

	fields := map[string]interface{}{
		UserIDKey:  FromUserIDContext(ctx),
		TraceIDKey: FromTraceIDContext(ctx),
		VersionKey: version,
	}
	if v := o.Title; v != "" {
		fields[SpanTitleKey] = v
	}
	if v := o.FuncName; v != "" {
		fields[SpanFunctionKey] = v
	}

	return newEntry(logrus.WithFields(fields))
}

// StartSpanWithCall ...
func StartSpanWithCall(ctx context.Context, opts ...SpanOption) func() *Entry {
	return func() *Entry {
		return StartSpan(ctx, opts...)
	}
}

// Debugf ...
func Debugf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Debugf(format, args...)
}

// Infof ...
func Infof(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Infof(format, args...)
}

// Printf ...
func Printf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Printf(format, args...)
}

// Warnf ...
func Warnf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Warnf(format, args...)
}

// Errorf ...
func Errorf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Errorf(format, args...)
}

// Fatalf ...
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Fatalf(format, args...)
}

func newEntry(entry *logrus.Entry) *Entry {
	return &Entry{entry: entry}
}

// Entry ...
type Entry struct {
	entry *logrus.Entry
}

func (e *Entry) checkAndDelete(fields map[string]interface{}, keys ...string) {
	for _, key := range keys {
		if _, ok := fields[key]; ok {
			delete(fields, key)
		}
	}
}

// WithFields ...
func (e *Entry) WithFields(fields map[string]interface{}) *Entry {
	e.checkAndDelete(fields,
		TraceIDKey,
		SpanTitleKey,
		SpanFunctionKey,
		VersionKey)
	return newEntry(e.entry.WithFields(fields))
}

// WithField ...
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(map[string]interface{}{key: value})
}

// Fatalf ...
func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.entry.Fatalf(format, args...)
}

// Errorf ...
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.entry.Errorf(format, args...)
}

// Warnf ...
func (e *Entry) Warnf(format string, args ...interface{}) {
	e.entry.Warnf(format, args...)
}

// Infof ...
func (e *Entry) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

// Printf ...
func (e *Entry) Printf(format string, args ...interface{}) {
	e.entry.Printf(format, args...)
}

// Debugf ...
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.entry.Debugf(format, args...)
}
