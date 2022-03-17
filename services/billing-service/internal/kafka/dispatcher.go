package kafka

import (
	"context"
	"encoding/json"

	"billing-service/internal/messaging"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Dispatcher struct {
	writer *kafka.Writer
}

func NewDispatcher(writer *kafka.Writer) *Dispatcher {
	return &Dispatcher{writer: writer}
}

func (d *Dispatcher) Dispatch(ctx context.Context, message messaging.Message) error {
	value, err := json.Marshal(message)
	if err != nil {
		return errors.Wrapf(err, `failed to marshal message "%s"`, message.Name())
	}

	err = d.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
		Headers: []kafka.Header{
			{Key: "name", Value: []byte(message.Name())},
		},
	})
	if err != nil {
		return errors.Wrapf(err, `failed to dispatch message "%s"`, message.Name())
	}

	return nil
}
