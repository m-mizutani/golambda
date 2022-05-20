package golambda

import (
	"errors"

	"github.com/m-mizutani/goerr"
	// use pkg/errors to generate stacktrace
)

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
