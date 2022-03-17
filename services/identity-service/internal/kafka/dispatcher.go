package kafka

import (
	"context"
	"encoding/json"

	"identity-service/internal/events"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Dispatcher struct {
	writer *kafka.Writer
}

func NewDispatcher(writer *kafka.Writer) *Dispatcher {
	return &Dispatcher{writer: writer}
}

func (d *Dispatcher) Dispatch(ctx context.Context, event events.Event) error {
	value, err := json.Marshal(event)
	if err != nil {
		return errors.Wrapf(err, `failed to marshal event "%s"`, event.Name())
	}

	err = d.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
		Headers: []kafka.Header{
			{Key: "name", Value: []byte(event.Name())},
		},
	})
	if err != nil {
		return errors.Wrapf(err, `failed to dispatch event "%s"`, event.Name())
	}

	return nil
}
