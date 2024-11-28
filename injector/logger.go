package injector

import (
	"os"
	"path/filepath"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
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

	return func() {
		if file != nil {
			file.Close()
		}
	}, nil
}
