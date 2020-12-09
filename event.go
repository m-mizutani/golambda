package golambda

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/m-mizutani/golambda/errors"
)

// Event provides lambda original event converting utilities
type Event struct {
	Origin interface{}
}

// Bind does json.Marshal original event and json.Unmarshal to v
func (x *Event) Bind(v interface{}) error {
	raw, err := json.Marshal(x.Origin)
	if err != nil {
		return errors.Wrap(err, "Marshal original lambda event").With("originalEvent", x.Origin)
	}

	if err := json.Unmarshal(raw, v); err != nil {
		return errors.Wrap(err, "Unmarshal to v").With("raw", string(raw))
	}

	return nil
}

// EventRecord is decapsulate event data (e.g. Body of SQS event)
type EventRecord []byte

// Bind unmarshal event record to object
func (x EventRecord) Bind(ev interface{}) error {
	if err := json.Unmarshal(x, ev); err != nil {
		return errors.Wrap(err, "Failed json.Unmarshal in DecodeEvent").With("raw", string(x))
	}
	return nil
}

// DecapSQSEvent decapsulate wrapped body data in SQSEvent
func (x *Event) DecapSQSEvent() ([]EventRecord, error) {
	var sqsEvent events.SQSEvent
	if err := x.Bind(&sqsEvent); err != nil {
		return nil, err
	}

	var output []EventRecord
	for _, record := range sqsEvent.Records {
		if record.MessageId == "" {
			continue
		}
		output = append(output, EventRecord(record.Body))
	}

	if len(output) == 0 {
		return nil, errors.New("No SQS event records")
	}

	return output, nil
}

// DecapSNSonSQSEvent decapsulate wrapped body data to in SNSEntity over SQSEvent
func (x *Event) DecapSNSonSQSEvent() ([]EventRecord, error) {
	var sqsEvent events.SQSEvent
	if err := x.Bind(&sqsEvent); err != nil {
		return nil, err
	}

	if len(sqsEvent.Records) == 0 {
		return nil, errors.New("No SQS event records")
	}

	var output []EventRecord
	for _, record := range sqsEvent.Records {
		var snsEntity events.SNSEntity
		if err := json.Unmarshal([]byte(record.Body), &snsEntity); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal SNS entity in SQS msg").With("body", record.Body)
		}

		output = append(output, EventRecord(snsEntity.Message))
	}

	return output, nil
}

// DecapSNSEvent decapsulate wrapped body data in SNSEvent
func (x *Event) DecapSNSEvent() ([]EventRecord, error) {
	var snsEvent events.SNSEvent
	if err := x.Bind(&snsEvent); err != nil {
		return nil, err
	}

	if len(snsEvent.Records) == 0 {
		return nil, errors.New("No SNS event records")
	}

	var output []EventRecord
	for _, record := range snsEvent.Records {
		output = append(output, EventRecord(record.SNS.Message))
	}

	return output, nil
}
