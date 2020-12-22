# golambda [![Travis-CI](https://travis-ci.com/m-mizutani/golambda.svg)](https://travis-ci.org/m-mizutani/golambda) [![Report card](https://goreportcard.com/badge/github.com/m-mizutani/golambda)](https://goreportcard.com/report/github.com/m-mizutani/golambda) [![Go Reference](https://pkg.go.dev/badge/github.com/m-mizutani/golambda.svg)](https://pkg.go.dev/github.com/m-mizutani/golambda)

A suite of Go utilities for AWS Lambda functions to ease adopting best practices.

## Overview
### Features

- **Event decapsulation**: Parse event data received when invoking. Also `golambda` make easy to write unit test of Lambda function
- **Structured logging**: `golambda` provides requisite minimum logging interface for Lambda function. It output log as structured JSON.
- **Error handling**: Error structure with arbitrary variables and stack trace feature.
- **Get parameters**: Secret values should be stored in AWS Secrets Manager and can be got easily.

NOTE: The suite is **NOT** focusing to Lambda function for API gateway, but partially can be leveraged for the function.

### How to use

```
$ go get github.com/m-mizutani/golambda
```


## Source event decapsulation

Lambda function can have event source(s) such as SQS, SNS, etc. The main data is encapsulated in their data structure. `golambda` provides not only decapsulation feature for Lambda execution but also encapsulation feature for testing. Following event sources are supported for now.

- SQS body: `DecapSQSBody`
- SNS message: `DecapSNSMessage`
- SNS message over SQS: `DecapSNSonSQSMessage`

### Lambda implementation

```go
package main

import (
	"strings"

	"github.com/m-mizutani/golambda"
)

// MyEvent is exported for test
type MyEvent struct {
	Message string `json:"message"`
}

// Handler is sample function. The function concatenates all message in SQS body and returns it.
func Handler(event golambda.Event) (interface{}, error) {
	// Decapsulate body message(s) in SQS Event structure.
	// It can also handles multiple SQS records.
	events, err := event.DecapSQSBody()
	if err != nil {
		return nil, err
	}

	var response []string

	// Iterate body message(S)
	for _, ev := range events {
		var msg MyEvent
		// Unmarshal golambda.Event to MyEvent
		if err := ev.Bind(&msg); err != nil {
			return nil, err
		}

		// Do something
		response = append(response, msg.Message)
	}

	return strings.Join(response, ":"), nil
}

func main() {
	golambda.Start(Handler)
}
```

### Unit test

```go
package main_test

import (
	"testing"

	"github.com/m-mizutani/golambda"
	"github.com/stretchr/testify/require"

	main "github.com/m-mizutani/golambda/example/decapEvent"
)

func TestHandler(t *testing.T) {
	var event golambda.Event
	messages := []main.MyEvent{
		{
			Message: "blue",
		},
		{
			Message: "orange",
		},
	}
	require.NoError(t, event.EncapSQS(messages))

	resp, err := main.Handler(event)
	require.NoError(t, err)
	require.Equal(t, "blue:orange", resp)
}
```

## Structured logging

Lambda function output log data to CloudWatch Logs by default. CloudWatch Logs and Insights that is rich CloudWatch Logs viewer supports JSON format logs. Therefore JSON formatted log is better for Lambda function.

`golambda` provides `Logger` for JSON format logging. It has `With()` and `Set()` to add a pair of key and value to a log message. `Logger` has lambda request ID by default if you use the logger with `golambda.Start()`.

### Output with temporary variable

```go
v1 := "say hello"
golambda.Logger.With("var1", v1).Info("Hello, hello, hello")
/* Output:
{
	"level": "info",
	"lambda.requestID": "565389dc-c13f-4fc0-b113-xxxxxxxxxxxx",
	"time": "2020-12-13T02:44:30Z",
	"var1": "say hello",
	"message": "Hello, hello, hello"
}
*/
```

### Set permanent variable to logger

```go
golambda.Logger.Set("myRequestID", myRequestID)

// ~~~~~~~ snip ~~~~~~

golambda.Logger.Error("oops")
/* Output:
{
	"level": "error",
	"lambda.requestID": "565389dc-c13f-4fc0-b113-xxxxxxxxxxxx",
	"time": "2020-11-12T02:44:30Z",
	"myRequestID": "xxxxxxxxxxxxxxxxx",
	"message": "oops"
}
*/
```

### Log level

`golambda.Logger` (`golambda.LambdaLogger` type) provides following log level. Log level can be configured by environment variable `LOG_LEVEL`.

- `TRACE`
- `DEBUG`
- `INFO`
- `ERROR`

Lambda function should return error to top level function when occurring unrecoverable error, should not exit suddenly. Therefore `PANIC` and `FATAL` is not provided according to the thought.

## Error handling

`golambda.Error` can have pairs of key and value to keep context of error. For example, `golambda.Error` can bring original string data when failed to unmarshal JSON. The string data can be extracted in caller function.

Also, `golambda.Start` supports general error handling:

1. Output error log with
    - Pairs of key and value in `golambda.Error` as `error.values`
    - Stack trace of error as `error.stacktrace`
2. Send error record to sentry.io if `SENTRY_DSN` is set as environment variable
    - Stack trace of `golambda.Error` is also available in sentry.io by compatibility with `github.com/pkg/errors`
    - Output event ID of sentry to log as `error.sentryEventID`

```go
package main

import (
	"github.com/m-mizutani/golambda"
)

// Handler is exported for test
func Handler(event golambda.Event) (interface{}, error) {
	trigger := "something wrong"
	return nil, golambda.NewError("oops").With("trigger", trigger)
}

func main() {
	golambda.Start(Handler)
}
```

Then, `golambda` output following log to CloudWatch.

```json
{
    "level": "error",
    "lambda.requestID": "565389dc-c13f-4fc0-b113-f903909dbd45",
    "trigger": "something wrong",
    "stacktrace": [
        {
            "func": "main.Handler",
            "file": "xxx/your/project/src/main.go",
            "line": 27
        },
        {
            "func": "github.com/m-mizutani/golambda.Start.func1",
            "file": "xxx/github.com/m-mizutani/golambda/lambda.go",
            "line": 107
        }
    ],
    "time": "2020-12-13T02:42:48Z",
    "message": "oops"
}
```

## License

See [LICENSE](LICENSE).
