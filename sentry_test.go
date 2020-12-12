package golambda

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSentry(t *testing.T) {
	if _, ok := os.LookupEnv("SENTRY_DSN"); !ok {
		t.Skip("SENTRY_DSN is not set")
	}

	id := emitSentry(NewError("oops"))
	assert.NotEmpty(t, id)
	flushSentry()
}
