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
