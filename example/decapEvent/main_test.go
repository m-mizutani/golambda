package main_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/golambda/v2"
	"github.com/stretchr/testify/require"

	main "github.com/m-mizutani/golambda/v2/example/decapEvent"
)

func TestHandler(t *testing.T) {
	messages := []main.MyEvent{
		{
			Message: "blue",
		},
		{
			Message: "orange",
		},
	}
	event, err := golambda.NewSQSEvent(messages)
	require.NoError(t, err)

	resp, err := main.Handler(context.Background(), event)
	require.NoError(t, err)
	require.Equal(t, "blue:orange", resp)
}
