package errors_test

import (
	"fmt"
	"testing"

	"github.com/m-mizutani/golambda/errors"
	"github.com/stretchr/testify/assert"
)

func oops() *errors.Error {
	return errors.New("omg")
}

func normalError() error {
	return fmt.Errorf("red")
}

func wrapError() *errors.Error {
	err := normalError()
	return errors.Wrap(err, "orange")
}

func TestNewError(t *testing.T) {
	err := oops()
	assert.Contains(t, fmt.Sprintf("%+v", err.Unwrap()), "errors_test.oops")
	assert.Contains(t, err.Error(), "omg")
}

func TestWrapError(t *testing.T) {
	err := wrapError()
	st := fmt.Sprintf("%+v", err.Unwrap())
	assert.Contains(t, st, "errors_test.wrapError")
	assert.NotContains(t, st, "errors_test.normalError")
	assert.Contains(t, err.Error(), "orange: red")
}

func TestSentryEmit(t *testing.T) {
	err := wrapError()
	errors.EmitSentry(err)
	errors.FlushSentry()
	_, ok := err.Values["sentry.eventID"]
	assert.False(t, ok) // SENTRY_DSN is not set
}
