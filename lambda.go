package golambda

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Callback is callback function type of golambda.Start().
//
// Trigger event data (SQS, SNS, etc) is included in Event.
//
// Callback has 2 returned value. 1st value (interface{}) will be passed to Lambda. The 1st value is allowed nil if you do not want to return any value to Lambda. 2nd value (error) also will be passed to Lambda, however golambda.Start() does error handling:
// 1) Extract stack trace of error if err is golambda.Error
// 2) Send error record to sentry.io if SENTRY_DSN is set as environment variable
// 3) Output error log
type Callback[T any] func(event *Event) (T, error)

// Start sets up Arguments and logging tools, then invoke Callback with Arguments. When exiting, it also does error handling if Callback returns error
func Start[T any](callback Callback[T]) {
	lambda.Start(func(ctx context.Context, origin interface{}) (interface{}, error) {
		defer flushSentry()

		lc, _ := lambdacontext.FromContext(ctx)
		initLogger()
		WithLogger("lambda.requestID", lc.AwsRequestID)

		Logger.With("input", origin).Info("Starting Lambda")

		event := NewEvent(ctx, origin)

		resp, err := callback(event)
		if err != nil {
			HandleError(err)
			return nil, err
		}

		Logger.With("output", resp).Info("Exiting Lambda")

		return resp, nil
	})
}
