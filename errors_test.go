package golambda_test

import (
	"fmt"
	"testing"

	"github.com/m-mizutani/golambda"
	"github.com/stretchr/testify/assert"
)

func oops() *golambda.Error {
	return golambda.NewError("omg")
}

func normalError() error {
	return fmt.Errorf("red")
}

func wrapError() *golambda.Error {
	err := normalError()
	return golambda.WrapError(err, "orange")
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
