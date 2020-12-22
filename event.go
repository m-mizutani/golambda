package golambda

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

// Event provides lambda original event converting utilities
type Event struct {
	Ctx    context.Context
	Origin interface{}
}

// Bind does json.Marshal original event and json.Unmarshal to v
func (x *Event) Bind(v interface{}) error {
	raw, err := json.Marshal(x.Origin)
	if err != nil {
		return WrapError(err, "Marshal original lambda event").With("originalEvent", x.Origin)
	}

	if err := json.Unmarshal(raw, v); err != nil {
		return WrapError(err, "Unmarshal to v").With("raw", string(raw))
	}

	return nil
}

// EventRecord is decapsulate event data (e.g. Body of SQS event)
type EventRecord []byte

// Bind unmarshal event record to object
func (x EventRecord) Bind(ev interface{}) error {
	if err := json.Unmarshal(x, ev); err != nil {
		return WrapError(err, "Failed json.Unmarshal in DecodeEvent").With("raw", string(x))
	}
	return nil
}

// String returns raw string data
func (x EventRecord) String() string {
	return string(x)
}

// DecapSQSBody decapsulate wrapped body data in SQSEvent
func (x *Event) DecapSQSBody() ([]EventRecord, error) {
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
		return nil, NewError("No SQS event records")
	}

	return output, nil
}

// EncapSQS sets v as SQSEvent body. This function overwrite Origin for testing.
// EncapSQS allows both of one record and multiple record as slice or array
// e.g.)
//   ev.EncapSQS("red") -> one SQSMessage in SQSEvent
//   ev.EncapSQS([]string{"blue", "orange"}) -> two SQSMessage in SQSEvent
func (x *Event) EncapSQS(v interface{}) error {
	messages, err := encapSQSMessage(v)
	if err != nil {
		return err
	}

	x.Origin = events.SQSEvent{
		Records: messages,
	}
	return nil
}

func encapSQSMessage(v interface{}) ([]events.SQSMessage, error) {
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
			return nil, WrapError(err, "Failed to marshal v")
		}
		return []events.SQSMessage{
			{
				MessageId: uuid.New().String(),
				Body:      string(raw),
			},
		}, nil
	}
}

// DecapSNSonSQSMessage decapsulate wrapped body data to in SNSEntity over SQSEvent
func (x *Event) DecapSNSonSQSMessage() ([]EventRecord, error) {
	var sqsEvent events.SQSEvent
	if err := x.Bind(&sqsEvent); err != nil {
		return nil, err
	}

	if len(sqsEvent.Records) == 0 {
		return nil, NewError("No SQS event records")
	}

	var output []EventRecord
	for _, record := range sqsEvent.Records {
		var snsEntity events.SNSEntity
		if err := json.Unmarshal([]byte(record.Body), &snsEntity); err != nil {
			return nil, WrapError(err, "Failed to unmarshal SNS entity in SQS msg").With("body", record.Body)
		}

		output = append(output, EventRecord(snsEntity.Message))
	}

	return output, nil
}

// EncapSNSonSQSMessage sets v as SNS entity over SQS. This function overwrite Origin and should be used for testing.
// EncapSNSonSQSMessage allows both of one record and multiple record as slice or array
//
// e.g.)
//
//     ev.EncapSNSonSQSMessage("red") // -> one SQS message on one SQS event
//     ev.EncapSNSonSQSMessage([]string{"blue", "orange"}) // -> two SQS message on one SQS event
func (x *Event) EncapSNSonSQSMessage(v interface{}) error {
	snsEntities, err := encapSNSEntity(v)
	if err != nil {
		return err
	}

	sqsMessages, err := encapSQSMessage(snsEntities)
	if err != nil {
		return err
	}

	x.Origin = events.SQSEvent{
		Records: sqsMessages,
	}

	return nil
}

// DecapSNSMessage decapsulate wrapped body data in SNSEvent
func (x *Event) DecapSNSMessage() ([]EventRecord, error) {
	var snsEvent events.SNSEvent
	if err := x.Bind(&snsEvent); err != nil {
		return nil, err
	}

	if len(snsEvent.Records) == 0 {
		return nil, NewError("No SNS event records")
	}

	var output []EventRecord
	for _, record := range snsEvent.Records {
		output = append(output, EventRecord(record.SNS.Message))
	}

	return output, nil
}

// EncapSNS sets v as SNS entity. This function overwrite Origin for testing.
// EncapSNS allows both of one record and multiple record as slice or array
// e.g.)
// ev.EncapSNS("red") -> one SNS entity in SNSEvent
// ev.EncapSNS([]string{"blue", "orange"}) -> two SNS entity in SNSEvent
func (x *Event) EncapSNS(v interface{}) error {
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

	x.Origin = events.SNSEvent{
		Records: snsRecords,
	}

	return nil
}

func encapSNSEntity(v interface{}) ([]events.SNSEntity, error) {
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
			return nil, WrapError(err, "Failed to marshal v")
		}
		return []events.SNSEntity{
			{
				Message: string(raw),
			},
		}, nil
	}
}
