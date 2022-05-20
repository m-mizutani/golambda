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

	resp, err := main.Handler(&event)
	require.NoError(t, err)
	require.Equal(t, "blue:orange", resp)
}
