package golambda

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Logger is common logging interface
var Logger zerolog.Logger

// InitLogger configures logger
func initLogger() {
	logLevel := os.Getenv("LOG_LEVEL")

	var zeroLogLevel zerolog.Level
	switch strings.ToLower(logLevel) {
	case "trace":
		zeroLogLevel = zerolog.TraceLevel
	case "debug":
		zeroLogLevel = zerolog.DebugLevel
	case "info":
		zeroLogLevel = zerolog.InfoLevel
	case "error":
		zeroLogLevel = zerolog.ErrorLevel
	default:
		zeroLogLevel = zerolog.InfoLevel
	}

	var writer io.Writer = zerolog.ConsoleWriter{Out: os.Stdout}
	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		// If running on AWS Lambda
		writer = os.Stdout
	}

	logger := zerolog.New(writer).Level(zeroLogLevel).With().Timestamp().Logger()
	Logger = logger
}

func setRequestIDtoLogger(requestID string) {
	Logger = Logger.With().Str("lambda.requestID", requestID).Logger()
}
