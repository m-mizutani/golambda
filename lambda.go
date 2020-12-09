package golambda

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/m-mizutani/golambda/errors"
	"github.com/m-mizutani/golambda/logging"
)

// Handler is callback function type of lambda.Run()
type Handler func(ctx context.Context, event Event) (interface{}, error)

// Start sets up Arguments and logging tools, then invoke handler with Arguments
func Start(handler Handler) {
	lambda.Start(func(ctx context.Context, origin interface{}) (interface{}, error) {
		defer errors.FlushSentry()
		logging.Logger.Info().Interface("event", origin).Msg("Lambda start")

		event := Event{
			Origin: origin,
		}

		response, err := handler(ctx, event)
		if err != nil {
			errors.EmitSentry(err)

			log := logging.Logger.Error()
			if e, ok := err.(*errors.Error); ok {
				for key, value := range e.Values {
					log = log.Interface(key, value)
				}
				log = log.Str("stacktrace", e.StackTrace())
			}

			log.Msg(err.Error())
			return nil, err
		}

		return response, nil
	})
}
