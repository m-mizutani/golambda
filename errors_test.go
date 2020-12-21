package golambda_test

import (
	"fmt"
	"testing"

	"github.com/m-mizutani/golambda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Contains(t, fmt.Sprintf("%+v", err), "golambda_test.oops")
	assert.Contains(t, err.Error(), "omg")
}

func TestWrapError(t *testing.T) {
	err := wrapError()
	st := fmt.Sprintf("%+v", err)
	assert.Contains(t, st, "github.com/m-mizutani/golambda_test.wrapError\n")
	assert.Contains(t, st, "github.com/m-mizutani/golambda_test.TestWrapError\n")
	assert.NotContains(t, st, "github.com/m-mizutani/golambda_test.normalError\n")
	assert.Contains(t, err.Error(), "orange: red")
}

func TestStackTrace(t *testing.T) {
	err := oops()
	st := err.Stacks()
	require.Equal(t, 4, len(st))
	assert.Equal(t, "github.com/m-mizutani/golambda_test.oops", st[0].Func)
	assert.Regexp(t, `/golambda/errors_test\.go$`, st[0].File)
	assert.Equal(t, 13, st[0].Line)
}
