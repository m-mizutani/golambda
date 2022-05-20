package golambda

import (
	"os"
	"strings"

	"github.com/m-mizutani/zlog"
)

// Logger is common logging interface
var Logger *zlog.Logger

func initLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(logLevel) {
	case "trace", "debug", "info", "warn", "error":
		// nothing to do
	default:
		logLevel = "info"
	}

	RenewLogger(
		zlog.WithLogLevel(logLevel),
		zlog.WithEmitter(
			zlog.NewJsonEmitter(),
		),
	)
}

func RenewLogger(options ...zlog.Option) {
	Logger = zlog.New(options...)
}

func WithLogger(key string, value interface{}) {
	Logger = Logger.With(key, value)
}
