package configs

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

var env EnvConfig

// EnvConfig ...
type EnvConfig struct {
	RunMode string `envconfig:"RUN_MODE" required:"true" default:"debug"` // (debug/test/release)

	// http
	HttpHost         string        `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	HttpPort         int           `envconfig:"HTTP_PORT" default:"8080"`
	HttpShutdownTime time.Duration `envconfig:"HTTP_SHUTDOWN_TIME" default:"30"`

	// log
	LogLevel      int    `envconfig:"LOG_LEVEL" default:"5"`       // 1:fatal 2:error,3:warn,4:info,5:debug
	LogFormat     string `envconfig:"LOG_FORMAT" default:"json"`   // (text/json)
	LogOutput     string `envconfig:"LOG_OUTPUT" default:"stdout"` // (stdout/stderr/file)
	LogOutputFile string `envconfig:"LOG_OUTPUT_FILE" default:""`  // etc: data/admin.log

	// gorm
	GormDebug             bool          `envconfig:"GORM_DEBUG" default:"true"`
	GormDbType            string        `envconfig:"GORM_DB_TYPE" default:"mysql"`
	GormMaxLifetime       time.Duration `envconfig:"GORM_MAX_LIFETIME" default:"7200s"`
	GormMaxOpenConns      int           `envconfig:"GORM_MAX_OPENCONNS" default:"150"`
	GormMaxIdleConns      int           `envconfig:"GORM_MAX_IDLE_CONNS" default:"50"`
	GormEnableAutoMigrate bool          `envconfig:"GORM_ENABLE_AUTO_MIGRATE" default:"true"`
	GormTablePrefix       string        `envconfig:"GORM_TABLE_PREFIX" default:""`

	// mysql
	MysqlHost       string `envconfig:"MYSQL_HOST" default:"127.0.0.1"`
	MysqlPort       int    `envconfig:"MYSQL_PORT" default:"3306"`
	MysqlUser       string `envconfig:"MYSQL_USER" default:"admin"`
	MysqlPassword   string `envconfig:"MYSQL_PASSWORD" default:"pass"`
	MysqlDbName     string `envconfig:"MYSQL_DB_NAME" default:"ddd-gin-admin"`
	MysqlParameters string `envconfig:"MYSQL_PARAMETERS" default:"charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true"`
}

// InitEnv init env.
func InitEnv() error {
	return envconfig.Process("", &env)
}

func Env() EnvConfig {
	return env
}
