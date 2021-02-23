package golambda

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// LambdaLogger provides basic logging features for Lambda function. golambda.Logger is configured by default as global variable of golambda.
type LambdaLogger struct {
	zeroLogger zerolog.Logger
}

// NewLambdaLogger returns a new LambdaLogger.
// NOTE: golambda.Logger is recommended for general usage.
func NewLambdaLogger(logLevel string) *LambdaLogger {
	var zeroLogLevel zerolog.Level
	switch strings.ToLower(logLevel) {
	case "trace":
		zeroLogLevel = zerolog.TraceLevel
	case "debug":
		zeroLogLevel = zerolog.DebugLevel
	case "info":
		zeroLogLevel = zerolog.InfoLevel
	case "warn":
		zeroLogLevel = zerolog.WarnLevel
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
	return &LambdaLogger{
		zeroLogger: logger,
	}
}

// LogEntry is one record of logging. Trace, Debug, Info and Error methods emit message and values
type LogEntry struct {
	logger *LambdaLogger
	values map[string]interface{}
}

// Entry returns a new LogEntry
func (x *LambdaLogger) Entry() *LogEntry {
	return &LogEntry{
		logger: x,
		values: make(map[string]interface{}),
	}
}

// Trace output log as Trace level message
func (x *LambdaLogger) Trace(msg string) { x.Entry().Trace(msg) }

// Debug output log as Debug level message
func (x *LambdaLogger) Debug(msg string) { x.Entry().Debug(msg) }

// Info output log as Info level message
func (x *LambdaLogger) Info(msg string) { x.Entry().Info(msg) }

// Warn output log as Warning level message
func (x *LambdaLogger) Warn(msg string) { x.Entry().Warn(msg) }

// Error output log as Error level message
func (x *LambdaLogger) Error(msg string) { x.Entry().Error(msg) }

// Set saves key and value to logger. The key and value are output permanently
func (x *LambdaLogger) Set(key string, value interface{}) {
	x.zeroLogger = x.zeroLogger.With().Interface(key, value).Logger()
}

// With adds key and value to log message. Value will be represented by zerolog.Interface
func (x *LambdaLogger) With(key string, value interface{}) *LogEntry {
	entry := x.Entry()
	entry.values[key] = value
	return entry
}

// With saves key and value into own and return own pointer.
func (x *LogEntry) With(key string, value interface{}) *LogEntry {
	x.values[key] = value
	return x
}

func (x *LogEntry) bind(ev *zerolog.Event) {
	for k, v := range x.values {
		ev.Interface(k, v)
	}
}

// Trace emits log message as trace level.
func (x *LogEntry) Trace(msg string) {
	ev := x.logger.zeroLogger.Trace()
	x.bind(ev)
	ev.Msg(msg)
}

// Debug emits log message as debug level.
func (x *LogEntry) Debug(msg string) {
	ev := x.logger.zeroLogger.Debug()
	x.bind(ev)
	ev.Msg(msg)
}

// Info emits log message as info level.
func (x *LogEntry) Info(msg string) {
	ev := x.logger.zeroLogger.Info()
	x.bind(ev)
	ev.Msg(msg)
}

// Warn emits log message as Warn level.
func (x *LogEntry) Warn(msg string) {
	ev := x.logger.zeroLogger.Warn()
	x.bind(ev)
	ev.Msg(msg)
}

// Error emits log message as error level.
func (x *LogEntry) Error(msg string) {
	ev := x.logger.zeroLogger.Error()
	x.bind(ev)
	ev.Msg(msg)
}

// Logger is common logging interface
var Logger *LambdaLogger

// InitLogger configures logger
func initLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	Logger = NewLambdaLogger(logLevel)
}
