package main

import (
	"context"
	"strings"

	"github.com/m-mizutani/golambda/v2"
)

// MyEvent is exported for test
type MyEvent struct {
	Message string `json:"message"`
}

// Handler is exported for test
func Handler(ctx context.Context, event *golambda.Event) (string, error) {
	// Decapsulate body message(s) in SQS Event structure
	events, err := event.DecapSQS()
	if err != nil {
		return "", err
	}

	var response []string
	// Iterate body message(S)
	for _, ev := range events {
		var msg MyEvent
		// Unmarshal golambda.Event to MyEvent
		if err := ev.Bind(&msg); err != nil {
			return "", err
		}

		// Do something
		response = append(response, msg.Message)
	}

	return strings.Join(response, ":"), nil
}

func main() {
	golambda.Start(Handler)
}
