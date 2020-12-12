package golambda

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type lambdaLogger struct {
	logger zerolog.Logger
}

// LogEntry is one record of logging. Trace, Debug, Info and Error methods emit message and values
type LogEntry struct {
	logger zerolog.Logger
	values map[string]interface{}
}

func (x *lambdaLogger) NewLogEntry() *LogEntry {
	return &LogEntry{
		values: make(map[string]interface{}),
	}
}

func (x *lambdaLogger) Trace(msg string) { x.NewLogEntry().Trace(msg) }
func (x *lambdaLogger) Debug(msg string) { x.NewLogEntry().Debug(msg) }
func (x *lambdaLogger) Info(msg string)  { x.NewLogEntry().Info(msg) }
func (x *lambdaLogger) Error(msg string) { x.NewLogEntry().Error(msg) }

func (x *lambdaLogger) Set(key string, value interface{}) {
	x.logger = x.logger.With().Interface(key, value).Logger()
}

func (x *lambdaLogger) With(key string, value interface{}) *LogEntry {
	entry := x.NewLogEntry()
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
	ev := x.logger.Trace()
	x.bind(ev)
	ev.Msg(msg)
}

// Debug emits log message as debug level.
func (x *LogEntry) Debug(msg string) {
	ev := x.logger.Debug()
	x.bind(ev)
	ev.Msg(msg)
}

// Info emits log message as info level.
func (x *LogEntry) Info(msg string) {
	ev := x.logger.Info()
	x.bind(ev)
	ev.Msg(msg)
}

// Error emits log message as error level.
func (x *LogEntry) Error(msg string) {
	ev := x.logger.Error()
	x.bind(ev)
	ev.Msg(msg)
}

// Logger is common logging interface
var Logger *lambdaLogger

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
	Logger = &lambdaLogger{
		logger: logger,
	}
}
