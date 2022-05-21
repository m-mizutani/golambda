package golambda

import (
	"encoding/json"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
)

// Event provides lambda original event converting utilities
type Event struct {
	origin any
}

// Bind does json.Marshal original event and json.Unmarshal to v. If failed, return error with ErrInvalidEventData.
func (x *Event) Bind(v any) error {
	raw, err := json.Marshal(x.origin)
	if err != nil {
		return ErrFailedDecodeEvent.Wrap(err).With("origin", x.origin)
	}

	if err := json.Unmarshal(raw, v); err != nil {
		return ErrFailedDecodeEvent.Wrap(err).With("marshaled", string(raw))
	}

	return nil
}

// NewEvent provides a new event with original data.
func NewEvent(ev any) *Event {
	return &Event{
		origin: ev,
	}
}

// EventRecord is decapsulate event data (e.g. Body of SQS event)
type EventRecord []byte

// Bind unmarshal event record to object. If failed, return error with ErrInvalidEventData.
func (x EventRecord) Bind(ev any) error {
	if err := json.Unmarshal(x, ev); err != nil {
		return ErrFailedDecodeEvent.Wrap(err).With("raw", string(x))
	}
	return nil
}

// String returns raw string data
func (x EventRecord) String() string {
	return string(x)
}

// DecapSQS decapsulate wrapped body data in SQSEvent. If no SQS records, it returns ErrNoEventData.
func (x *Event) DecapSQS() ([]EventRecord, error) {
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
		return nil, goerr.Wrap(ErrNoEventData, "no SQS event records")
	}

	return output, nil
}

// DecapSNS decapsulate wrapped body data in SNSEvent. If no SQS records, it returns ErrNoEventData.
func (x *Event) DecapSNS() ([]EventRecord, error) {
	var snsEvent events.SNSEvent
	if err := x.Bind(&snsEvent); err != nil {
		return nil, err
	}

	if len(snsEvent.Records) == 0 {
		return nil, goerr.Wrap(ErrNoEventData, "no SNS event records")
	}

	var output []EventRecord
	for _, record := range snsEvent.Records {
		output = append(output, EventRecord(record.SNS.Message))
	}

	return output, nil
}

// DecapSNSoverSQS decapsulate wrapped body data to in SNSEntity over SQSEvent. If no SQS records, it returns ErrNoEventData.
func (x *Event) DecapSNSoverSQS() ([]EventRecord, error) {
	var sqsEvent events.SQSEvent
	if err := x.Bind(&sqsEvent); err != nil {
		return nil, err
	}

	if len(sqsEvent.Records) == 0 {
		return nil, goerr.Wrap(ErrNoEventData, "no SQS event records")
	}

	var output []EventRecord
	for _, record := range sqsEvent.Records {
		var snsEntity events.SNSEntity
		if err := json.Unmarshal([]byte(record.Body), &snsEntity); err != nil {
			return nil, goerr.Wrap(err, "Failed to unmarshal SNS entity in SQS msg").With("body", record.Body)
		}

		output = append(output, EventRecord(snsEntity.Message))
	}

	if len(output) == 0 {
		return nil, goerr.Wrap(ErrNoEventData, "no SNS event records")
	}

	return output, nil
}

// NewSQSEvent sets v as SQSEvent body. This function overwrite Origin for testing.
// NewSQSEvent allows both of one record and multiple record as slice or array
// e.g.)
//   ev.NewSQSEvent("red") -> one SQSMessage in SQSEvent
//   ev.NewSQSEvent([]string{"blue", "orange"}) -> two SQSMessage in SQSEvent
func NewSQSEvent(v any) (*Event, error) {
	messages, err := encapSQSMessage(v)
	if err != nil {
		return nil, err
	}

	return &Event{
		origin: events.SQSEvent{
			Records: messages,
		},
	}, nil
}

func encapSQSMessage(v any) ([]events.SQSMessage, error) {
	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		var messages []events.SQSMessage

		for i := 0; i < value.Len(); i++ {
			msg, err := encapSQSMessage(value.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg...)
		}

		return messages, nil

	case reflect.Ptr, reflect.UnsafePointer:
		if value.IsZero() || value.Elem().IsZero() {
			return nil, nil
		}

		return encapSQSMessage(value.Elem().Interface())

	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, goerr.Wrap(err, "Failed to marshal v")
		}
		return []events.SQSMessage{
			{
				MessageId: uuid.New().String(),
				Body:      string(raw),
			},
		}, nil
	}
}

// NewSNSonSQSEvent sets v as SNS entity over SQS. This function overwrite Origin and should be used for testing.
// NewSNSonSQSEvent allows both of one record and multiple record as slice or array
//
// e.g.)
//
//     ev.NewSNSonSQSEvent("red") // -> one SQS message on one SQS event
//     ev.NewSNSonSQSEvent([]string{"blue", "orange"}) // -> two SQS message on one SQS event
func (x *Event) NewSNSonSQSEvent(v any) error {
	snsEntities, err := encapSNSEntity(v)
	if err != nil {
		return err
	}

	sqsMessages, err := encapSQSMessage(snsEntities)
	if err != nil {
		return err
	}

	x.origin = events.SQSEvent{
		Records: sqsMessages,
	}

	return nil
}

// NewSNSEvent sets v as SNS entity. This function overwrite Origin for testing.
// NewSNSEvent allows both of one record and multiple record as slice or array
// e.g.)
// ev.NewSNSEvent("red") -> one SNS entity in SNSEvent
// ev.NewSNSEvent([]string{"blue", "orange"}) -> two SNS entity in SNSEvent
func (x *Event) NewSNSEvent(v any) error {
	entities, err := encapSNSEntity(v)
	if err != nil {
		return err
	}

	var snsRecords []events.SNSEventRecord
	for _, entity := range entities {
		snsRecords = append(snsRecords, events.SNSEventRecord{
			SNS: entity,
		})
	}

	x.origin = events.SNSEvent{
		Records: snsRecords,
	}

	return nil
}

func encapSNSEntity(v any) ([]events.SNSEntity, error) {
	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		var messages []events.SNSEntity

		for i := 0; i < value.Len(); i++ {
			msg, err := encapSNSEntity(value.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg...)
		}

		return messages, nil

	case reflect.Ptr, reflect.UnsafePointer:
		if value.IsZero() || value.Elem().IsZero() {
			return nil, nil
		}

		return encapSNSEntity(value.Elem().Interface())

	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, goerr.Wrap(err, "Failed to marshal v")
		}
		return []events.SNSEntity{
			{
				Message: string(raw),
			},
		}, nil
	}
}
