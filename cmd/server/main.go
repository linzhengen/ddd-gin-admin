package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/s"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/logger"
)

// VERSION app version，
// can also be set via：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "5.2.1"

var (
	swaggerDir string
)

func init() {
	flag.StringVar(&swaggerDir, "swagger", "", "swagger dir")
}

func main() {
	flag.Parse()

	var state int32 = 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger.SetVersion(VERSION)
	logger.SetTraceIDFunc(s.NewTraceID)
	ctx := logger.NewTraceIDContext(context.Background(), s.NewTraceID())
	span := logger.StartSpanWithCall(ctx)

	call := Init(ctx)
EXIT:
	for {
		sig := <-sc
		span().Printf("catch signal [%s]", sig.String())

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			atomic.StoreInt32(&state, 0)
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	if call != nil {
		call()
	}

	span().Printf("stopped server")
	time.Sleep(time.Second)
	os.Exit(int(atomic.LoadInt32(&state)))
}
