package main

import (
	"os"
	"path/filepath"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/logger"
)

// InitLogger init logger
func InitLogger() (func(), error) {
	c := configs.Env()
	logger.SetLevel(c.LogLevel)
	logger.SetFormatter(c.LogFormat)

	var file *os.File
	if c.LogOutput != "" {
		switch c.LogOutput {
		case "stdout":
			logger.SetOutput(os.Stdout)
		case "stderr":
			logger.SetOutput(os.Stderr)
		case "file":
			if name := c.LogOutputFile; name != "" {
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

	return func() {
		if file != nil {
			file.Close()
		}
	}, nil
}
