package injector

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
	loggerhook "github.com/linzhengen/ddd-gin-admin/pkg/logger/hook"
	loggergormhook "github.com/linzhengen/ddd-gin-admin/pkg/logger/hook/gorm"
	"github.com/sirupsen/logrus"
)

func InitLogger() (func(), error) {
	c := configs.C.Log
	logger.SetLevel(c.Level)
	logger.SetFormatter(c.Format)

	// log output
	var file *os.File
	if c.Output != "" {
		switch c.Output {
		case "stdout":
			logger.SetOutput(os.Stdout)
		case "stderr":
			logger.SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
				_ = os.MkdirAll(filepath.Dir(name), 0777)

				f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					return nil, err
				}
				logger.SetOutput(f)
				file = f
			}
		}
	}

	var hook *loggerhook.Hook
	if c.EnableHook {
		var hookLevels []logrus.Level
		for _, lvl := range c.HookLevels {
			plvl, err := logrus.ParseLevel(lvl)
			if err != nil {
				return nil, err
			}
			hookLevels = append(hookLevels, plvl)
		}

		if c.Hook.IsGorm() {
			hc := configs.C.LogGormHook

			var dsn string
			switch hc.DBType {
			case "mysql":
				dsn = configs.C.MySQL.DSN()
			case "sqlite3":
				dsn = configs.C.Sqlite3.DSN()
			case "postgres":
				dsn = configs.C.Postgres.DSN()
			default:
				return nil, errors.New("unknown db")
			}

			h := loggerhook.New(loggergormhook.New(&loggergormhook.Config{
				DBType:       hc.DBType,
				DSN:          dsn,
				MaxLifetime:  hc.MaxLifetime,
				MaxOpenConns: hc.MaxOpenConns,
				MaxIdleConns: hc.MaxIdleConns,
				TableName:    hc.Table,
			}),
				loggerhook.SetMaxWorkers(c.HookMaxThread),
				loggerhook.SetMaxQueues(c.HookMaxBuffer),
				loggerhook.SetLevels(hookLevels...),
			)
			logger.AddHook(h)
			hook = h
		}
	}

	return func() {
		if file != nil {
			file.Close()
		}

		if hook != nil {
			hook.Flush()
		}
	}, nil
}
