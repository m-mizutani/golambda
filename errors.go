package golambda

import (
	"errors"

	"github.com/m-mizutani/goerr"
	// use pkg/errors to generate stacktrace
)

// HandleError emits the error to Sentry and outputs the error to logs
func HandleError(err error) {
	logger := Logger

	if evID := emitSentry(err); evID != "" {
		logger = logger.With("error.sentryEventID", evID)
	}

	var goErr *goerr.Error
	if errors.As(err, &goErr) {
		logger = logger.With("error.values", goErr.Values())
		logger = logger.With("error.stacktrace", goErr.Stacks())
	}

	logger.Error(err.Error())
}

var (
	ErrFailedDecodeEvent = goerr.New("failed to decode event")
	ErrFailedEncodeEvent = goerr.New("failed to encode event")

	ErrNoEventData = goerr.New("no event data")

	ErrInvalidARN           = goerr.New("invalid ARN")
	ErrFailedSecretsManager = goerr.New("failed SecretsManager operation")
	ErrFailedDecodeSecret   = goerr.New("failed to decode secret")
)
