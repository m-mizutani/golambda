# golambda [![Travis-CI](https://travis-ci.com/m-mizutani/golambda.svg)](https://travis-ci.org/m-mizutani/golambda) [![GoDoc](https://godoc.org/github.com/m-mizutani/golambda?status.svg)](http://godoc.org/github.com/m-mizutani/golambda) [![Report card](https://goreportcard.com/badge/github.com/m-mizutani/golambda)](https://goreportcard.com/report/github.com/m-mizutani/golambda)

Utilities for lambda function in Go language.

Use cases

- Decapsulate source event
- Logging
- Error handling (with sentry if needed)

## Decapsulate source event

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