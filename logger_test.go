package golambda

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {

	t.Run("output trace level or more when set trace level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := &LambdaLogger{
			zeroLogger: zerolog.New(buf).Level(zerolog.TraceLevel),
		}

		logger.Trace("out/Trace")
		assert.Contains(t, buf.String(), "out/Trace")
		logger.Debug("out/Debug")
		assert.Contains(t, buf.String(), "out/Debug")
		logger.Info("out/Info")
		assert.Contains(t, buf.String(), "out/Info")
		logger.Error("out/Error")
		assert.Contains(t, buf.String(), "out/Error")
	})

	t.Run("output debug level or more when set debug level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := &LambdaLogger{
			zeroLogger: zerolog.New(buf).Level(zerolog.DebugLevel),
		}

		logger.Trace("out/Trace")
		assert.NotContains(t, buf.String(), "out/Trace")
		logger.Debug("out/Debug")
		assert.Contains(t, buf.String(), "out/Debug")
		logger.Info("out/Info")
		assert.Contains(t, buf.String(), "out/Info")
		logger.Error("out/Error")
		assert.Contains(t, buf.String(), "out/Error")
	})

	t.Run("output info level or more when set info level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := &LambdaLogger{
			zeroLogger: zerolog.New(buf).Level(zerolog.InfoLevel),
		}

		logger.Trace("out/Trace")
		assert.NotContains(t, buf.String(), "out/Trace")
		logger.Debug("out/Debug")
		assert.NotContains(t, buf.String(), "out/Debug")
		logger.Info("out/Info")
		assert.Contains(t, buf.String(), "out/Info")
		logger.Error("out/Error")
		assert.Contains(t, buf.String(), "out/Error")
	})

	t.Run("output error level or more when set error level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := &LambdaLogger{
			zeroLogger: zerolog.New(buf).Level(zerolog.ErrorLevel),
		}

		logger.Trace("out/Trace")
		assert.NotContains(t, buf.String(), "out/Trace")
		logger.Debug("out/Debug")
		assert.NotContains(t, buf.String(), "out/Debug")
		logger.Info("out/Info")
		assert.NotContains(t, buf.String(), "out/Info")
		logger.Error("out/Error")
		assert.Contains(t, buf.String(), "out/Error")
	})

	t.Run("output values provided by With", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := &LambdaLogger{
			zeroLogger: zerolog.New(buf).Level(zerolog.InfoLevel),
		}

		logger.With("color", "blue").With("magic", "five").Info("not sane")
		assert.Contains(t, buf.String(), `"color":"blue"`)
		assert.Contains(t, buf.String(), `"magic":"five"`)
		assert.Contains(t, buf.String(), `"not sane"`)
	})

	t.Run("output values provided by Set", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := &LambdaLogger{
			zeroLogger: zerolog.New(buf).Level(zerolog.InfoLevel),
		}

		logger.Set("color", "orange")
		logger.Info("dolls")

		assert.Contains(t, buf.String(), `"color":"orange"`)
		assert.Contains(t, buf.String(), `"dolls"`)
	})
}
