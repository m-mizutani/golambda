package logging_test

import (
	"github.com/m-mizutani/golambda/logging"
)

func ExampleLogger() {
	logging.Logger.Info().Msg("hoge")
}
