package golambda

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Callback is callback function type of lambda.Run()
type Callback func(event Event) error

// Start sets up Arguments and logging tools, then invoke Callback with Arguments
func Start(callback Callback) {
	lambda.Start(func(ctx context.Context, origin interface{}) error {
		defer flushSentry()
		Logger.Info().Interface("event", origin).Msg("Lambda start")

		lc, _ := lambdacontext.FromContext(ctx)
		setRequestIDtoLogger(lc.AwsRequestID)

		event := Event{
			Ctx:    ctx,
			Origin: origin,
		}

		if err := callback(event); err != nil {
			log := Logger.Error()

			if evID := emitSentry(err); evID != "" {
				log = log.Str("sentry.eventID", evID)
			}

			if e, ok := err.(*Error); ok {
				for key, value := range e.Values() {
					log = log.Interface(key, value)
				}
				log = log.Interface("stacktrace", e.StackTrace())
			}

			log.Msg(err.Error())
			return err
		}

		return nil
	})
}
