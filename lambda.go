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
type Callback func(event Event) (interface{}, error)

// Start sets up Arguments and logging tools, then invoke Callback with Arguments. When exiting, it also does error handling if Callback returns error
func Start(callback Callback) {
	lambda.Start(func(ctx context.Context, origin interface{}) (interface{}, error) {
		defer flushSentry()

		initLogger()

		lc, _ := lambdacontext.FromContext(ctx)
		Logger.Set("lambda.requestID", lc.AwsRequestID)

		Logger.With("event", origin).Info("Lambda start")

		event := Event{
			Ctx:    ctx,
			Origin: origin,
		}

		resp, err := callback(event)
		if err != nil {
			EmitError(err)
			return nil, err
		}

		return resp, nil
	})
}
