# golambda [![Travis-CI](https://travis-ci.com/m-mizutani/golambda.svg)](https://travis-ci.org/m-mizutani/golambda) [![GoDoc](https://godoc.org/github.com/m-mizutani/golambda?status.svg)](http://godoc.org/github.com/m-mizutani/golambda) [![Report card](https://goreportcard.com/badge/github.com/m-mizutani/golambda)](https://goreportcard.com/report/github.com/m-mizutani/golambda)

Utilities for lambda function in Go language.

Use cases

- Decapsulate source event
- Logging
- Error handling (with sentry)

## Decapsulate source event

```go
package main

import (
	"context"
	"fmt"

	"github.com/m-mizutani/golambda"
)

type myMessage struct {
	Something string
}

func Handler(ctx context.Context, event golambda.Event) (interface{}, error) {
	events, err := event.DecapSQSBody()
	if err != nil {
		return nil, err
	}

	for _, ev := range events {
		var msg myMessage
		err := ev.Bind(&msg)
		if err != nil {
			return nil, err
		}

		// Do something
		fmt.Println("something: ", msg.Something)
	}

	return nil, nil
}

func main() {
	golambda.Start(Handler)
}
```
