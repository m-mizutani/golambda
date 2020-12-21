# golambda [![Travis-CI](https://travis-ci.com/m-mizutani/golambda.svg)](https://travis-ci.org/m-mizutani/golambda) [![Report card](https://goreportcard.com/badge/github.com/m-mizutani/golambda)](https://goreportcard.com/report/github.com/m-mizutani/golambda) [![Go Reference](https://pkg.go.dev/badge/github.com/m-mizutani/golambda.svg)](https://pkg.go.dev/github.com/m-mizutani/golambda)

A suite of Go utilities for AWS Lambda functions to ease adopting best practices.

### Features

- **Source event decapsulation**: Parse event data received when invoking. Also `golambda` make easy to write unit test of Lambda function
- **Structured logging**: `golambda` provides requisite minimum logging interface for Lambda function. It output log as structured JSON.
- **Error handling**: Error structure with arbitrary variables and stack trace feature.

## Source event decapsulation

- SQS body
- SNS message
- SNS message over SQS
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

// Handler is exported for test
func Handler(event golambda.Event) (interface{}, error) {
	// Decapsulate body message(s) in SQS Event structure
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

When occurring unrecoverable error, it should return as error to top level function in Lambda. Therefore `PANIC` and `FATAL` is not provided.

## Error handling

`golambda.Error` can have pairs of key and value to know context of error. In `golambda.Start`, error log with key and value is output to CloudWatch as structured log when returned error is `golambda.Error`. Example is below.

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
            "line": 27
        },
		// -------- snip --------------
    ],
    "time": "2020-12-13T02:42:48Z",
    "message": "oops"
}
```

## License

See [LICENSE](LICENSE).
