package golambda

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Callback is callback function type of golambda.Start()
type Callback func(event Event) (interface{}, error)

// Start sets up Arguments and logging tools, then invoke Callback with Arguments
func Start(callback Callback) {
	lambda.Start(func(ctx context.Context, origin interface{}) (interface{}, error) {
		defer flushSentry()
		Logger.With("event", origin).Info("Lambda start")

		lc, _ := lambdacontext.FromContext(ctx)
		Logger.Set("lambda.requestID", lc.AwsRequestID)

		event := Event{
			Ctx:    ctx,
			Origin: origin,
		}

		resp, err := callback(event)
		if err != nil {
			entry := Logger.Entry()

			if evID := emitSentry(err); evID != "" {
				entry.With("sentry.eventID", evID)
			}

			if e, ok := err.(*Error); ok {
				for key, value := range e.Values() {
					entry = entry.With(key, value)
				}
				entry.With("stacktrace", e.StackTrace())
			}

			entry.Error(err.Error())
			return nil, err
		}

		return resp, nil
	})
}
