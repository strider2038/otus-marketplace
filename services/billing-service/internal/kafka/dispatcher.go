package kafka

import (
	"context"
	"encoding/json"

	"billing-service/internal/messaging"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type Dispatcher struct {
	writer *kafka.Writer
	logger zerolog.Logger
}

func NewDispatcher(writer *kafka.Writer, logger zerolog.Logger) *Dispatcher {
	return &Dispatcher{writer: writer, logger: logger}
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

	d.logger.Info().
		Str("name", message.Name()).
		RawJSON("body", value).
		Msg("Message dispatched")

	return nil
}
