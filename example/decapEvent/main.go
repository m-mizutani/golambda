// +build decap_event

package main

import (
	"fmt"

	"github.com/m-mizutani/golambda"
)

type myMessage struct {
	Something string
}

func Handler(event golambda.Event) (interface{}, error) {
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
