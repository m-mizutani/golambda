package golambda_test

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-mizutani/golambda"
)

func TestDecapSQSEvent(t *testing.T) {
	t.Run("can make SQSEvent to EventRecord", func(t *testing.T) {
		v := &golambda.Event{
			Origin: events.SQSEvent{
				Records: []events.SQSMessage{
					{
						MessageId: "t1",
						Body:      "blue",
					},
					{
						MessageId: "t2",
						Body:      "orange",
					},
				},
			},
		}
		events, err := v.DecapSQSEvent()

		require.NoError(t, err)
		require.Equal(t, 2, len(events))
		assert.Equal(t, "blue", string(events[0]))
		assert.Equal(t, "orange", string(events[1]))
	})

	t.Run("fail when no SQS event", func(t *testing.T) {
		v := &golambda.Event{
			Origin: events.SQSEvent{},
		}
		events, err := v.DecapSQSEvent()

		require.Error(t, err)
		assert.Nil(t, events)
	})

	t.Run("fail with SNS event", func(t *testing.T) {
		v := &golambda.Event{
			Origin: events.SNSEvent{
				Records: []events.SNSEventRecord{
					{
						SNS: events.SNSEntity{
							Message: "blue",
						},
					},
					{
						SNS: events.SNSEntity{
							Message: "orange",
						},
					},
				},
			},
		}
		events, err := v.DecapSQSEvent()

		require.Error(t, err)
		assert.Nil(t, events)
	})
}
